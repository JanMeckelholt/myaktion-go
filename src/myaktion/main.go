package main

import (
	"log"
	"net/http"

	"github.com/JanMeckelholt/myaktion-go/src/myaktion/handler"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting My-Aktion API server")
	router := mux.NewRouter()
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.HandleFunc("/campaign", handler.CreateCampaign).Methods("POST")
	router.HandleFunc("/campaigns", handler.GetCampaigns).Methods("GET")
	if err := http.ListenAndServe(":8002", router); err != nil {
		log.Fatal(err)
	}

}
