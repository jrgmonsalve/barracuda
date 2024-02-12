package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jrgmonsalve/barracuda/src/zincsearch"
	"github.com/rs/cors"
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
	r.Use(corsMiddleware)
	r.Post("/emails/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request: ", r)

		var req searchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		result, err := zincsearch.QueryByEmailField(req.FieldName, req.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Println(result)
		w.Write(result)
	})
	r.Get("/emails/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log.Println("Fetching email with ID: ", id)

		result, err := zincsearch.QueryByEmailId(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if result == nil {
			http.Error(w, "Email not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Println(result)
		w.Write(result)
	})

	http.ListenAndServe(":"+port, r)
}

func corsMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"}, // Ajusta esta lista según tu caso
		AllowCredentials: true,
		Debug:            true,
	}).Handler(next)
}
