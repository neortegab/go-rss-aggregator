package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func loadDotEnv() {
	data, err := os.ReadFile(".env")

	if err != nil {
		log.Fatal("error:", errors.New("there was an error opening .env file"))
	}

	fileData := string(data)
	pairVals := strings.Split(fileData, "\n")

	for _, v := range pairVals {
		pair := strings.Split(v, "=")
		key, val := pair[0], pair[1]
		os.Setenv(key, val)
	}
	log.Print(".env loaded")
}

func main() {
	loadDotEnv()
	port := os.Getenv("PORT")

	if port == "" {
		err := errors.New("error: PORT not found")
		log.Fatalf("%v\n", err)
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	nRouter := chi.NewRouter()
	nRouter.Get("/healthz", handlerReadiness)
	nRouter.Get("/err", handlerErr)

	router.Mount("/v1", nRouter)

	log.Printf("Initializing on server port %s...\n", port)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalf("Error at serving\n")
		return
	}
}
