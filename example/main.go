package main

import (
	"log"
	"net/http"

	docs "github.com/depozzyx/mux-docs"
	"github.com/gorilla/mux"
)

var host string = "localhost:8000"

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/sendData", endpoint).Methods("POST", "PUT")
	router.HandleFunc("/addConsumer", endpoint).Methods("PUT").Name("Adds and returns id")
	router.HandleFunc("/editConsumer/{consumerID:.*}", endpoint).Methods("PATCH").Queries("name", "{name}")
	router.HandleFunc("/deleteConsumer/{consumerID:.*}", endpoint).Methods("DELETE")
	router.HandleFunc("/consumeDataIDs/{consumerID:.*}", endpoint).Methods("GET")
	router.HandleFunc("/consumeData/{consumerID}/{dataID}", endpoint).Methods("GET")

	log.Printf("Listening on %s\n", host)

	docsMiddleware := docs.Middleware(router, host)
	log.Fatal(http.ListenAndServe(host, docsMiddleware(router)))
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}
