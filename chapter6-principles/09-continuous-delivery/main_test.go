package main

import (
	"os"
	"testing"
)

func TestFeatureFlagParse(t *testing.T) {
	os.Setenv("FEATURE_SPLIT_BILL", "true")
	flag := os.Getenv("FEATURE_SPLIT_BILL") == "true"
	if !flag {
		t.Errorf("Expected feature flag to be true")
	}

	os.Setenv("FEATURE_SPLIT_BILL", "false")
	flag = os.Getenv("FEATURE_SPLIT_BILL") == "true"
	if flag {
		t.Errorf("Expected feature flag to be false")
	}
}