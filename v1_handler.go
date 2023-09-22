package main

import (
	"net/http"

	v1handlers "github.com/AtinAgnihotri/glog-aggregator/v1_handlers"
	"github.com/go-chi/chi/v5"
)

func V1Handler(serverConf *ServerConf) http.Handler {
	r := chi.NewRouter()

	r.Get("/readiness", v1handlers.Readiness)

	r.Get("/err", v1handlers.Err)

	return r
}
