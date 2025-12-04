package util

import (
	"math"
	"testing"
)

func TestDistanceKm_SamePointIsZero(t *testing.T) {
	d := DistanceKm(41.0, 29.0, 41.0, 29.0)
	if d > 0.001 {
		t.Errorf("same point distance should be ~0, got %f", d)
	}
}

func TestDistanceKm_IstanbulToAnkara(t *testing.T) {
	// İstanbul: 41.0082, 28.9784
	// Ankara:   39.9334, 32.8597
	d := DistanceKm(41.0082, 28.9784, 39.9334, 32.8597)

	if math.Abs(d-351) > 10 {
		t.Errorf("istanbul-ankara mesafesi ~351km olmalı, got %f", d)
	}
}
