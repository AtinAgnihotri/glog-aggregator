package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ServerConf struct {
	PORT string
	DB   *database.Queries
}

func main() {

	serverConf := ServerConf{}

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env variables")
	}

	dbURL := os.Getenv("DB_URL")

	if len(dbURL) <= 0 {
		log.Fatal("Unable to fetch DB URL")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Unable to connect to db")
	}

	dbQueries := database.New(db)

	serverConf.DB = dbQueries
	serverConf.PORT = ":" + os.Getenv("PORT")

	if len(serverConf.PORT) == 1 {
		serverConf.PORT = ""
	}

	router := chi.NewRouter()
	corsHandler := cors.Handler(cors.Options{})
	router.Use(corsHandler)

	// Mount v1 Handlers
	router.Mount("/v1", V1Handler(&serverConf))

	server := &http.Server{
		Addr:    serverConf.PORT,
		Handler: router,
	}

	log.Printf("Listening on port %v", serverConf.PORT)
	log.Fatal(server.ListenAndServe())
}
