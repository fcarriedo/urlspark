package store

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
func (db *memStore) Persist(longURL string, ttl int) (string, error) {
	db.Lock()
	defer db.Unlock()

	for {
		id := genRandID(idLen)

		if _, exists := db.m[id]; !exists {
			// If not existent SET it with with the expiration window
			db.m[id] = longURL
			// launch the expiration routine
			time.AfterFunc(time.Duration(ttl)*time.Second, func() {
				db.Del(id)
			})

			return id, nil
		}
	}

	return "", fmt.Errorf("Could not store %s", longURL)
}

func (db *memStore) Get(id string) (string, error) {
	db.RLock()
	defer db.RUnlock()

	url, exists := db.m[id]
	if !exists {
		return "", fmt.Errorf("entry does not exist or has expired")
	}

	return url, nil
}

func (db *memStore) Del(id string) error {
	db.Lock()
	defer db.Unlock()

	delete(db.m, id)

	return nil
}
