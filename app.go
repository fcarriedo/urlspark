package main

import (
	"flag"
	"fmt"
	"github.com/fcarriedo/urlspark/store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	port      = flag.Int("p", 80, "http port to run")
	ttl       = flag.Int("ttl", 60, "The expiration time [time to live seconds]")
	redisAddr = flag.String("redis", "", "redis address 'host:port'")
)

// Our URL datastore
var ds store.UrlStore

func init() {
	if *redisAddr != "" {
		// If redis param given, init the datastore with it.
		var err error
		if ds, err = store.NewRedisStore(*redisAddr); err != nil {
			// Fail fast if no datastore at startup
			log.Fatal(err)
		}
	} else {
		ds = store.NewMemoryStore()
	}

	// Flags usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:  %v [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Starts a URL shortener server\n\n")
		flag.PrintDefaults()
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

		id, err := ds.Persist(url, *ttl)
		if err != nil {
			log.Println(err)
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
		url, err := ds.Get(id)
		if err != nil {
			http.Error(w, "The requested URL has expired and/or does not exist", http.StatusNotFound)
			return
		}
		http.Redirect(w, req, url, http.StatusFound)
	case "DELETE":
		ds.Del(id)
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
