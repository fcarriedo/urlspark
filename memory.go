package main

import (
	"fmt"
	"sync"
	"time"
)

type memStore struct {
	sync.RWMutex
	m map[string]string
}

func NewMemoryStore() *memStore {
	return &memStore{m: make(map[string]string)}
}

// Persists the given URL and returns the unique ID that references it
func (db *memStore) persist(longURL string, expSec int) (string, error) {
	db.Lock()
	defer db.Unlock()

	for {
		id := genRandID(idLen)

		if _, exists := db.m[id]; !exists {
			// If not existent SET it with with the expiration window
			db.m[id] = longURL
			// launch the expiration routine
			go func() {
				<-time.After(time.Duration(expSec) * time.Second)
				db.del(id)
			}()

			return id, nil
		}
	}

	return "", fmt.Errorf("Could not store %s", longURL)
}

func (db *memStore) get(id string) (string, error) {
	db.RLock()
	defer db.RUnlock()

	url, exists := db.m[id]
	if !exists {
		return "", fmt.Errorf("entry does not exist or has expired")
	}

	return url, nil
}

func (db *memStore) del(id string) error {
	db.Lock()
	defer db.Unlock()

	delete(db.m, id)

	return nil
}
