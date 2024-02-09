package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jrgmonsalve/barracuda/src/zincsearch"
)

type searchRequest struct {
	FieldName string `json:"fieldName"`
	Value     string `json:"value"`
}

func main() {

	var port string

	flag.StringVar(&port, "port", "3000", "define el puerto en el cual el servidor estará escuchando")
	flag.Parse()
	log.Println("Starting server on port: ", port)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/emails/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request: ", r)

		var req searchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		result, err := zincsearch.QueryByEmail(req.FieldName, req.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		w.Write(result)
	})

	http.ListenAndServe(":"+port, r)
}
