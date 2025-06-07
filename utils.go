package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// DNS request result
type Result struct {
	Duration time.Duration
	Failed   bool
}

// Get string env
func getEnv(key string, defaultVal string) string {
	if val, exist := os.LookupEnv(key); exist {
		return val
	}
	return defaultVal
}

// Get env as Duration
func getEnvAsDuration(key string, defaultVal int) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(value) * time.Millisecond
	}

	return time.Duration(defaultVal) * time.Millisecond
}

// Request to DNS server
func dnsRequest(wg *sync.WaitGroup, dns_server string, client *dns.Client, fqdn string, c chan<- Result) {
	defer wg.Done()
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	_, rtt, err := client.Exchange(m, dns_server+":53")
	if err != nil {
		slog.Error(err.Error())
		c <- Result{rtt, true}
	}
	c <- Result{rtt, false}
}

// Calculate time
func calculateTime(c <-chan Result) string {
	var min time.Duration
	var max time.Duration
	var totalTime time.Duration
	var failedAmount int
	for val := range c {

		if val.Failed {
			failedAmount++
			continue
		}

		if min == 0 || val.Duration < min {
			min = val.Duration
		}
		if val.Duration > max {
			max = val.Duration
		}
		totalTime += val.Duration
	}

	avgTime := totalTime / time.Duration(cap(c)-failedAmount)

	if cap(c) == 0 || avgTime == 0 {
		return "No successful requests"
	}

	return fmt.Sprintf("Requests amount: %d, Fastest response time: %s. Slowest response time: %s Average response time: %s \n", cap(c), min, max, avgTime)
}

// DNS stress func
func stressTest(count int, delay time.Duration, dns_server string, dns_list []string) string {
	var wg sync.WaitGroup
	wg.Add(count)
	c := make(chan Result, count)
	dns_client := &dns.Client{}

	randRange := func() int {
		return rand.Intn(len(dns_list))
	}

	limiter := time.Tick(delay)
	for range count {
		<-limiter
		go dnsRequest(&wg, dns_server, dns_client, dns_list[randRange()], c)
	}

	wg.Wait()
	close(c)

	return calculateTime(c)
}
