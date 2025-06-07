package main

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	dnsServer      string
	fqdnListFile   string
	delay          time.Duration
	requestsAmount []int
}

func newConfig() *Config {
	return &Config{
		dnsServer:      getEnv("DNS_SERVER", "8.8.8.8"),
		fqdnListFile:   getEnv("FQDN_LIST_PATH", "./dns_list.txt"),
		delay:          getEnvAsDuration("DELAY", 2),
		requestsAmount: getEnvAsIntSlice("REQUESTS_AMOUNT", []int{10}),
	}
}

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	config := newConfig()

	slog.Info("Starting up.", "Using DNS server", config.dnsServer)

	err := checkIPAvailability(config.dnsServer, "53")
	if err != nil {
		slog.Error("DNS server is unavaliabe", "Addr", config.dnsServer)
		os.Exit(1)
	}
	slog.Info("DNS server is ok")

	file, err := os.Open(config.fqdnListFile)
	if err != nil {
		slog.Error("Cannot open file with dns list to resolve.", "Error:", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("Error closing file", "err", err)
		}
	}()

	var dns_list []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		dns_list = append(dns_list, scanner.Text())
	}

	slog.Info("File with fqdn list parsed successfully.", "Tolal amount of fqdn's", len(dns_list))

	slog.Info("Starting stress tests...")

	for _, count := range config.requestsAmount {
		_, err := io.WriteString(os.Stdout, stressTest(count, config.delay, config.dnsServer, dns_list))
		if err != nil {
			slog.Error("Cannot write answer", "Error:", err.Error())
		}
	}

	slog.Info("Test end's. Restart container to repeat. Untill that i'll sleep")

	for {
		time.Sleep(time.Hour * 24)
	}

}
