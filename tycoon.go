package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	mu                 sync.Mutex
	resources          int
	money              int
	resourceMultiplier int = 1
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/tycoon.html")
	})

	r.HandleFunc("/getTotals", getTotals).Methods("GET")
	r.HandleFunc("/getResourceMultiplier", getResourceMultiplierHandler).Methods("GET")
	r.HandleFunc("/makeResources", makeResourcesHandler).Methods("POST")
	r.HandleFunc("/sellResources", sellResourcesHandler).Methods("POST")
	r.HandleFunc("/2xCreateResources", twoxCreateResourceHandler).Methods("POST")

	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error: ", err)
	}
}

func getTotals(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Resources int `json:"resources"`
		Money     int `json:"money"`
	}{
		Resources: resources,
		Money:     money,
	})
}

func getResourceMultiplierHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		ResourceMultiplier int `json:"resourceMultiplier"`
	}{
		ResourceMultiplier: resourceMultiplier,
	})
}

func makeResourcesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	resources = resources + resourceMultiplier

	response := struct {
		Resources int `json:"resources"`
		Money     int `json:"money"`
	}{
		Resources: resources,
		Money:     money,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sellResourcesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	money += resources
	resources = 0

	response := struct {
		Money     int `json:"money"`
		Resources int `json:"resources"`
	}{
		Money:     money,
		Resources: resources,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func twoxCreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if money < 50 {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := struct {
			Error string `json:"error"`
		}{
			Error: "Not enough money",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	money -= 50
	resourceMultiplier = 2

	response := struct {
		Money     int `json:"money"`
		Resources int `json:"resources"`
	}{
		Money:     money,
		Resources: resources,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
