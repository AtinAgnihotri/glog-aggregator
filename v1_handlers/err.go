package v1handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
)

func Err(w http.ResponseWriter, _ *http.Request) {
	err := helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to respond on /err: %v", err.Error()))
	}
}
