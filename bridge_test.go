package main

import (
	"math"
	"testing"

	"github.com/tikv/client-go/v2/oracle"
)

func TestGCSafePointFromCurrentTS(t *testing.T) {
	currentTS := oracle.ComposeTS(120_000, 42)

	safePoint, err := gcSafePointFromCurrentTS(currentTS, 60)
	if err != nil {
		t.Fatalf("gcSafePointFromCurrentTS returned error: %v", err)
	}

	expected := oracle.ComposeTS(60_000, 0)
	if safePoint != expected {
		t.Fatalf("safePoint = %d, want %d", safePoint, expected)
	}
}

func TestGCSafePointFromCurrentTSRejectsZeroLifetime(t *testing.T) {
	if _, err := gcSafePointFromCurrentTS(oracle.ComposeTS(120_000, 0), 0); err == nil {
		t.Fatal("expected error for zero GC lifetime")
	}
}

func TestGCSafePointFromCurrentTSRejectsTooLargeLifetime(t *testing.T) {
	if _, err := gcSafePointFromCurrentTS(oracle.ComposeTS(120_000, 0), 121); err == nil {
		t.Fatal("expected error for GC lifetime larger than current physical time")
	}
}

func TestGCSafePointFromCurrentTSRejectsOverflowLifetime(t *testing.T) {
	overflowSeconds := uint64(math.MaxUint64/1000 + 1)

	if _, err := gcSafePointFromCurrentTS(oracle.ComposeTS(120_000, 0), overflowSeconds); err == nil {
		t.Fatal("expected error for overflowing GC lifetime")
	}
}
