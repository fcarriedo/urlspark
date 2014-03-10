package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var (
	port           = flag.Int("p", 80, "http port to run")
	expirationSecs = flag.Int("exp", 60, "The expiration time [seconds]")
	redisAddr      = flag.String("redis", "localhost:6379", "redis address 'host:port'")
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

func creationHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Fprintf(w, "URL shortener service")
	case "POST":
		url := req.FormValue("url")
		if url == "" {
			http.Error(w, "Required 'url' parameter missing", http.StatusBadRequest)
			return
		}

		id, err := persistURL(url)
		if err != nil {
			http.Error(w, "The server experienced an error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "http://%s/%s\r\n", req.Host, id)
	default:
		handleErr(w, http.StatusMethodNotAllowed)
	}
}

// Handles the mapping entry details.
func redirectHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	switch req.Method {
	case "GET":
		url, err := getURL(id)
		if err != nil {
			http.Error(w, "The requested URL has expired and/or does not exist", http.StatusNotFound)
			return
		}
		http.Redirect(w, req, url, http.StatusFound)
	case "DELETE":
		deleteURL(id)
	default:
		handleErr(w, http.StatusMethodNotAllowed)
	}
}

// Formats the given status in a standard. Any status would be managed in
// the same way whether is an error or not (does not enforce)
func handleErr(w http.ResponseWriter, errStatus int) {
	err := fmt.Sprintf("%d %s", errStatus, http.StatusText(errStatus))
	http.Error(w, err, errStatus)
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

func main() {
	flag.Parse()

	// The mux router
	router := mux.NewRouter()

	// Base URL handler
	router.HandleFunc("/", creationHandler)
	// Redirect handler
	router.HandleFunc("/{id}", redirectHandler)

	// Hook it with http pkg
	http.Handle("/", router)

	host := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server up and listening on %s\n", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
