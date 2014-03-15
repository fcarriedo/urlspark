package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisStore struct {
	pool *redis.Pool
}

// Creates a new URL datastore
func NewRedisStore(addr string) (*redisStore, error) {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				err := fmt.Errorf("redis is unreachable: %s", err)
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	s := &redisStore{pool: pool}

	// Test basic connectivity
	if err := s.ping(); err != nil {
		return nil, err
	}

	return s, nil
}

// Persists the given URL and returns the unique ID that references it
func (s *redisStore) persist(longURL string, ttl int) (string, error) {
	for {
		conn := s.pool.Get()
		defer conn.Close()

		id := genRandID(idLen)

		if exists, _ := redis.Bool(conn.Do("EXISTS", id)); !exists {
			// If not existent in redis, SET it with with the expiration window
			conn.Send("MULTI")
			conn.Send("SET", id, longURL)
			conn.Send("EXPIRE", id, ttl)
			if _, err := conn.Do("EXEC"); err != nil {
				return "", err
			}

			return id, nil
		}
	}
}

func (s *redisStore) get(id string) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", id))
}

func (s *redisStore) del(id string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", id)
	return err
}

// Do PING. Useful for connectivity and isAlive?
func (s *redisStore) ping() error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err
}
