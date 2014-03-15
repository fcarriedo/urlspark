package store

import "crypto/rand"

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const idLen = 4

// The interface definition for the URL datastore
type UrlStore interface {

	// Persist the given URL for the given amount of sec and returns the stored
	// URL identifier
	Persist(longUrl string, ttl int) (string, error)

	// Gets the stored URL given the identifier
	Get(id string) (string, error)

	// Deletes the URL given the identifier
	Del(id string) error
}

// Generates a random ID of the given length
func genRandID(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
