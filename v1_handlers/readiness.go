package v1handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
)

func Readiness(w http.ResponseWriter, _ *http.Request) {
	err := helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to respond on readiness: %v", err.Error()))
	}
}
