package main

import (
	"testing"
	"time"
)

func TestPersist(t *testing.T) {
	store := NewMemoryStore()
	longUrl := "http://very.com/long/url/hopefully/more/than/this"

	totalIds := 10
	for i := 0; i < totalIds; i++ {
		if _, err := store.persist(longUrl, 1); err != nil {
			t.Errorf("Error persisting: %s", err)
		}
	}

	if len(store.m) != 10 {
		t.Errorf("There should be 10 elements in the store")
	}

	// Check all generated random IDS are of the required length
	for k, _ := range store.m {
		if len(k) != idLen {
			t.Errorf("'%s' key is not of required length %d", k, idLen)
		}
	}

	// Lets wait until TTL expires (bit more than 1 sec)
	// Note: Sucks to wait on a test but helps us verify the ephemeral
	// nature of the store.
	time.Sleep(1100 * time.Millisecond)

	if len(store.m) != 0 {
		t.Errorf("The store should be empty by now. It has %s elems", len(store.m))
	}
}

func TestGet(t *testing.T) {
	store := NewMemoryStore()
	longUrl := "http://very.com/long/url/hopefully/more/than/this"

	id, _ := store.persist(longUrl, 5)

	url, _ := store.get(id)

	if url != longUrl {
		t.Errorf("Error getting the store URL through its ID")
	}
}

func TestDelete(t *testing.T) {
	store := NewMemoryStore()
	longUrl := "http://very.com/long/url/hopefully/more/than/this"

	id, _ := store.persist(longUrl, 5)

	store.del(id)

	if url, _ := store.get(id); url != "" {
		t.Errorf("The URL should not exist after deletion.")
	}
}
