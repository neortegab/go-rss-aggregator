package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"www.github.com/neortegab/go-rss-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

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

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL not found")
	}

	conn, errDb := sql.Open("postgres", dbURL)

	if errDb != nil {
		log.Fatalf("Can't connect to database, error: %v\n", errDb)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
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
	nRouter.Post("/users", apiCfg.handlerCreateUser)

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
