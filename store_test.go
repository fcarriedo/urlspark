package main

import "testing"

func TestRandIdGeneration(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := genRandID(i)
		if len(id) != i {
			t.Errorf("Expected ID of length %d but was of %d on %s", i, len(id), id)
		}
	}
}
