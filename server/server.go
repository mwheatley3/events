package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var (
	mu        = &sync.Mutex{}
	allEvents = []Event{}
)

// possible types for events
const (
	Increment = "INCREMENT"
	Decrement = "DECREMENT"
	BaseURL   = "http://localhost"
	Port      = "8080"
)

// Event is the structure for events
type Event struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

// Start starts the main api server
func Start() {
	r := router()
	fmt.Printf("Starting server on the port %s...", Port)

	log.Fatal(http.ListenAndServe(":"+Port, r))
}

func router() *mux.Router {

	router := mux.NewRouter()

	buildDir := "../../js/app/build/"
	// API routes
	router.HandleFunc("/event", increment).Methods("POST")
	router.HandleFunc("/events", all).Methods("GET")

	router.HandleFunc("/value/{t}", valueAt).Methods("GET")
	router.HandleFunc("/value", value).Methods("GET")

	// React routes
	buildHandler := http.FileServer(http.Dir(buildDir))
	router.PathPrefix("/").Handler(buildHandler)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(buildDir)))
	router.PathPrefix("/static/").Handler(staticHandler)

	return router
}

func valueAt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	println("\n valueAt route \n")
	fmt.Printf("\nparams %#+v\n", params)

	t, err := strconv.Atoi(params["t"])
	if err != nil || t < 0 {
		http.Error(w, fmt.Sprintf("Invalid Params: %#+v", params), http.StatusBadRequest)
		return
	}

	v := calcSystemValue(t)
	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(v)
	w.Write(js)
}

func value(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	v := calcSystemValue(-1)
	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(v)
	w.Write(js)
}

func all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	js, _ := json.Marshal(allEvents)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func increment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var e Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if e.Type != Increment && e.Type != Decrement {
		http.Error(w, fmt.Sprintf("unexpected type [%s] encountered", e.Type), http.StatusBadRequest)
		return
	}

	mu.Lock()
	allEvents = append(allEvents, e)
	mu.Unlock()

	js, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func calcSystemValue(endVal int) int {
	v := 0
	events := allEvents
	if endVal >= 0 {
		println("\n Updating all events \n")
		events = events[:endVal]
	}
	fmt.Printf("allEvents %#+v", allEvents)
	for _, e := range events {
		if e.Type == Increment {
			v += e.Value
		} else if e.Type == Decrement {
			v -= e.Value
		} else {
			fmt.Printf("unexpected type [%s] encountered", e.Type)
		}
	}
	return v
}
