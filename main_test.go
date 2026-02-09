package main

import "testing"

func TestRun_MissingArgs(t *testing.T) {
	err := run([]string{"cmd"})
	if err == nil {
		t.Fatal("expected error when args are missing")
	}
}

func TestRun_Success(t *testing.T) {
	err := run([]string{"cmd", "hello world"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
