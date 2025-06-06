package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV"
	expected := "value"
	err := os.Setenv(key, expected)

	if err != nil {
		fmt.Printf("Error with set env: %v", err)
	}

	defer func() {
		err := os.Unsetenv(key)
		if err != nil {
			fmt.Printf("Error with set env: %v", err)
		}
	}()

	if val := getEnv(key, "default"); val != expected {
		t.Errorf("Expected %s, got %s", expected, val)
	}

	if val := getEnv("UNDEFINED_ENV", "default"); val != "default" {
		t.Errorf("Expected default fallback, got %s", val)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	err := os.Setenv("INT_ENV", "42")
	if err != nil {
		fmt.Printf("Error with set env: %v", err)
	}

	defer func() {
		err := os.Unsetenv("INT_ENV")
		if err != nil {
			fmt.Printf("Error with set env: %v", err)
		}
	}()

	val := getEnvAsInt("INT_ENV", 10)
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	val = getEnvAsInt("MISSING_ENV", 99)
	if val != 99 {
		t.Errorf("Expected 99, got %d", val)
	}

	err = os.Setenv("INVALID_INT", "abc")
	if err != nil {
		fmt.Printf("Error with set env: %v", err)
	}

	val = getEnvAsInt("INVALID_INT", 55)
	if val != 55 {
		t.Errorf("Expected 55 fallback, got %d", val)
	}
}

func TestCalculateTime_AllSuccess(t *testing.T) {
	c := make(chan Result, 3)
	c <- Result{Duration: 10 * time.Millisecond, Failed: false}
	c <- Result{Duration: 20 * time.Millisecond, Failed: false}
	c <- Result{Duration: 30 * time.Millisecond, Failed: false}
	close(c)

	got := calculateTime(c)
	if got == "No successful requests" {
		t.Error("Expected results, got no successful requests")
	}

	if !contains(got, "Fastest response time") {
		t.Errorf("Expected proper output, got %s", got)
	}
}

func TestCalculateTime_WithFailures(t *testing.T) {
	c := make(chan Result, 3)
	c <- Result{Duration: 10 * time.Millisecond, Failed: false}
	c <- Result{Duration: 0, Failed: true}
	c <- Result{Duration: 20 * time.Millisecond, Failed: false}
	close(c)

	got := calculateTime(c)
	if !contains(got, "Average response time: 15ms") {
		t.Errorf("Expected real average, got: %s", got)
	}
}

func TestCalculateTime_NoSuccess(t *testing.T) {
	c := make(chan Result, 3)
	c <- Result{Duration: 0, Failed: true}
	c <- Result{Duration: 0, Failed: true}
	c <- Result{Duration: 0, Failed: true}
	close(c)

	got := calculateTime(c)
	if !contains(got, "No successful requests") {
		t.Errorf("Expected no successes, got: %s", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s != "" && s != "\n") && (stringContains(s, substr))
}

func stringContains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && s[:len(sub)] != "" && stringIndex(s, sub) >= 0)
}

func stringIndex(s, substr string) int {
	return len(s[:]) - len(([]rune(s))[len([]rune(s))-len(substr):])
}
