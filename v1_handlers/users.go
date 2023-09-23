package v1handlers

import (
	"net/http"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/google/uuid"
)

type UserRequest struct {
	Name string `json:"name"`
}

func (v1 *V1Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userRequest := UserRequest{}
	decodeErr := helpers.DecodeBodyToJson[UserRequest](r, &userRequest)

	if decodeErr != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Incorrect request recieved")
		return
	}

	usr, err := v1.DB.CreateUser(r.Context(), database.CreateUserParams{
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

func (v1 *V1Handlers) GetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	helpers.RespondWithJSON(w, http.StatusOK, user)
}
