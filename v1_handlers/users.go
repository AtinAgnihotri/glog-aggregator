package v1handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/auth"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/google/uuid"
)

type UserRequest struct {
	Name string `json:"name"`
}

type UserHandlers struct {
	DB *database.Queries
}

func (usrHandler *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userRequest := UserRequest{}
	decodeErr := helpers.DecodeBodyToJson[UserRequest](r, &userRequest)

	if decodeErr != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Incorrect request recieved")
		return
	}

	usr, err := usrHandler.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userRequest.Name,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, usr)

}

func (usrHandler *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey, err := auth.GetAuthApiKey(r)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "ApiKey not found")
		return
	}

	usr, err := usrHandler.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("User with ApiKey %v Not Found", apiKey))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, usr)
}
