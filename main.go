package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type ServerConf struct {
	PORT string
}

func V1Handler(serverConf *ServerConf) http.Handler {
	r := chi.NewRouter()

	return r
}

func main() {

	serverConf := ServerConf{}

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env variables")
	}

	serverConf.PORT = ":" + os.Getenv("PORT")

	if len(serverConf.PORT) == 1 {
		serverConf.PORT = ""
	}

	router := chi.NewRouter()
	corsHandler := cors.Handler(cors.Options{})
	router.Use(corsHandler)

	router.Mount("/v1", V1Handler(&serverConf))

	server := &http.Server{
		Addr:    serverConf.PORT,
		Handler: router,
	}

	log.Printf("Listening on port %v", serverConf.PORT)
	log.Fatal(server.ListenAndServe())
}
