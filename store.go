package main

import (
	"crypto/rand"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", *redisAddr)
			if err != nil {
				log.Fatalf("redis is unreachabe at '%s'", *redisAddr)
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

// Persists the given URL and returns the unique ID that references it
func persistURL(longURL string) (string, error) {
	for {
		conn := pool.Get()
		defer conn.Close()

		id := genRandID(4)

		if exists, _ := redis.Bool(conn.Do("EXISTS", id)); !exists {
			// If not existent in redis, SET it with with the expiration window
			conn.Send("MULTI")
			conn.Send("SET", id, longURL)
			conn.Send("EXPIRE", id, *expirationSecs)
			if _, err := conn.Do("EXEC"); err != nil {
				return "", err
			}

			return id, nil
		}
	}
}

// Gets the URL associated with the given ID
func getURL(id string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", id))
}

// Deletes the URL associated with the given ID
func deleteURL(id string) {
	conn := pool.Get()
	defer conn.Close()

	conn.Do("DEL", id)
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
