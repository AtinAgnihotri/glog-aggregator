package v1handlers

import (
	"fmt"
	"net/http"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/auth"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
)

type V1Handlers struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (v1 *V1Handlers) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		apiKey, err := auth.GetAuthApiKey(r)

		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, "ApiKey not found")
			return
		}

		usr, err := v1.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("User with ApiKey %v Not Found", apiKey))
			return
		}

		handler(w, r, usr)
	}
}
