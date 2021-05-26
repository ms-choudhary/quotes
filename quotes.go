package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ms-choudhary/quotes/store"
)

var (
	dbConnection = flag.String("db", "postgres://postgres@inst-juwcpfvrxydjkgptsanbzhlm-postgres-srv.dep-ns-inst-juwcpfvrxydjkgptsanbzhlm/production?sslmode=disable", "postgres db connection string")
)

var db store.DB

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func handleErr(w http.ResponseWriter, status int, errorMsg string, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
	log.Printf("%v: %v", errorMsg, err)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "failed to read request", err)
		return
	}

	var q store.Quote
	if err := json.Unmarshal(body, &q); err != nil {
		handleErr(w, http.StatusBadRequest, "couldn't parse body", err)
		return
	}

	q, err = db.Create(q)
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "couldn't store quote", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(q); err != nil {
		log.Printf("couldn't write body: %v", err)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := db.GetAll()
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "couldn't read quotes", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		log.Printf("couldn't write body: %v", err)
		return
	}
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	randQuote, err := db.GetRandom()
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "couldn't read quotes", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(randQuote); err != nil {
		log.Printf("couldn't write body: %v", err)
		return
	}
}

func main() {
	flag.Parse()

	var err error
	log.Printf("init db %v...", *dbConnection)
	db, err = store.NewPostgresStore(*dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Print("db connected.")

	r := mux.NewRouter()
	r.HandleFunc("/", healthcheckHandler)
	r.HandleFunc("/quotes/new", createHandler).Methods("POST")
	r.HandleFunc("/quotes/", indexHandler).Methods("GET")
	r.HandleFunc("/quotes/random", randomHandler).Methods("GET")

	http.Handle("/", r)

	log.Print("starting server...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
