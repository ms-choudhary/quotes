package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Quote struct {
	ID     int
	Text   string
	Author string
	Tags   []string
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request", http.StatusBadRequest)
	}

	var q Quote
	if err := json.Unmarshal(body, &q); err != nil {
		http.Error(w, "couldn't parse body", http.StatusUnprocessableEntity)
		return
	}

	q, err = Create(q)
	if err != nil {
		http.Error(w, "couldn't store quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(q); err != nil {
		http.Error(w, "couldn't write body", http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := GetAll()
	if err != nil {
		http.Error(w, "couldn't read quotes", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		http.Error(w, "couldn't write body", http.StatusInternalServerError)
	}
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	randQuote, err := GetRandom()
	if err != nil {
		http.Error(w, "couldn't read quotes", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(randQuote); err != nil {
		http.Error(w, "couldn't write body", http.StatusInternalServerError)
	}
}

func main() {
	err := Init("postgres", "postgres://postgres@localhost/production?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", healthcheckHandler)
	r.HandleFunc("/quotes/new", createHandler).Methods("POST")
	r.HandleFunc("/quotes/", indexHandler).Methods("GET")
	r.HandleFunc("/quotes/random", randomHandler).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9090", nil))
}