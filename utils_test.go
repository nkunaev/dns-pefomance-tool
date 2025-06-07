package main

import (
	"fmt"
	"os"
	"reflect"
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
			fmt.Printf("Error with unset env: %v", err)
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
			fmt.Printf("Error with unset env: %v", err)
		}
	}()

	val := getEnvAsDuration("INT_ENV", 10)
	if val != time.Duration(42)*time.Millisecond {
		t.Errorf("Expected 42, got %d", val)
	}

	val = getEnvAsDuration("MISSING_ENV", 99)
	if val != time.Duration(99)*time.Millisecond {
		t.Errorf("Expected 99, got %d", val)
	}

	err = os.Setenv("INVALID_INT", "abc")
	if err != nil {
		fmt.Printf("Error with set env: %v", err)
	}

	val = getEnvAsDuration("INVALID_INT", 55)
	if val != time.Duration(55)*time.Millisecond {
		t.Errorf("Expected 55 fallback, got %d", val)
	}
}

func TestGetEnvAsIntSlice(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal []int
		expected   []int
	}{
		{
			name:       "Valid integers",
			envKey:     "INT_SLICE_TEST_1",
			envValue:   "1,2,3,4",
			defaultVal: []int{9, 9, 9},
			expected:   []int{1, 2, 3, 4},
		},
		{
			name:       "Mixed valid and invalid integers",
			envKey:     "INT_SLICE_TEST_2",
			envValue:   "10,abc,20,xyz",
			defaultVal: []int{5, 5},
			expected:   []int{10, 20},
		},
		{
			name:       "All invalid integers",
			envKey:     "INT_SLICE_TEST_3",
			envValue:   "abc,def,ghi",
			defaultVal: []int{100, 200},
			expected:   []int{100, 200},
		},
		{
			name:       "Empty environment variable",
			envKey:     "INT_SLICE_TEST_4",
			envValue:   "",
			defaultVal: []int{7, 8, 9},
			expected:   []int{7, 8, 9},
		},
		{
			name:       "Whitespace and extra commas",
			envKey:     "INT_SLICE_TEST_5",
			envValue:   "1,,2,   ,3",
			defaultVal: []int{0},
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(tt.envKey, tt.envValue)
			if err != nil {
				fmt.Printf("Error with set env: %v", err)
			}

			result := getEnvAsIntSlice(tt.envKey, tt.defaultVal)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}

			err = os.Unsetenv(tt.envKey)
			if err != nil {
				fmt.Printf("Error with unset env: %v", err)
			}
		})
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
