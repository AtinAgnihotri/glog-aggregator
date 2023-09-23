package main

import (
	"net/http"

	v1handlers "github.com/AtinAgnihotri/glog-aggregator/v1_handlers"
	"github.com/go-chi/chi/v5"
)

func V1Handler(serverConf *ServerConf) http.Handler {
	r := chi.NewRouter()

	v1Handlers := v1handlers.V1Handlers{
		DB: serverConf.DB,
	}

	r.Get("/readiness", v1handlers.Readiness)

	r.Get("/err", v1handlers.Err)

	r.Post("/users", v1Handlers.CreateUser)

	r.Get("/users", v1Handlers.MiddlewareAuth(v1Handlers.GetUser))

	r.Post("/feeds", v1Handlers.MiddlewareAuth(v1Handlers.CreateFeed))

	return r
}
