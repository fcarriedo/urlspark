package store

import (
	"testing"
	"time"
)

const urlPrefix = "http://very.com/long/url/hopefully/more/than/this/"

func TestPersist(t *testing.T) {
	store := NewMemoryStore()

	totalEntries := 1000
	for i := 0; i < totalEntries; i++ {
		longUrl := urlPrefix + string(i)
		if _, err := store.Persist(longUrl, 1); err != nil {
			t.Errorf("%s (error) while persisting: %s", err, longUrl)
		}
	}

	if len(store.m) != totalEntries {
		t.Errorf("There should be %d elements in the store", totalEntries)
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
		t.Errorf("The store should be empty by now but it has %d elems", len(store.m))
	}
}

func TestGet(t *testing.T) {
	store := NewMemoryStore()

	for i := 0; i < 25; i++ {
		longUrl := urlPrefix + string(i)

		id, _ := store.Persist(longUrl, 5)

		url, _ := store.Get(id)

		if url != longUrl {
			t.Errorf("Expected '%s' but got '%s' using ID: %s", longUrl, url, id)
		}
	}
}

func TestDelete(t *testing.T) {
	store := NewMemoryStore()

	longUrl := urlPrefix + "1"
	id, _ := store.Persist(longUrl, 5)

	store.Del(id)

	if url, _ := store.Get(id); url != "" {
		t.Errorf("The URL should not exist after deletion.")
	}
}

// Benchmark tests

func BenchmarkPersist(b *testing.B) {
	store := NewMemoryStore()

	longUrl := urlPrefix + "1"

	for i := 0; i < b.N; i++ {
		store.Persist(longUrl, 1)
	}
}

func BenchmarkGet(b *testing.B) {
	store := NewMemoryStore()

	longUrl := urlPrefix + "1"
	id, _ := store.Persist(longUrl, 200)

	for i := 0; i < b.N; i++ {
		store.Get(id)
	}
}
