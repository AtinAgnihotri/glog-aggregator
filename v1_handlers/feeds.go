package v1handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/google/uuid"
)

type FeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// func (v1 *V1Handlers) GetAllFeeds(w http.ResponseWriter, r *http.Request) {}

func (v1 *V1Handlers) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedReq := FeedRequest{}
	err := helpers.DecodeBodyToJson[FeedRequest](r, &feedReq)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request recived. Body must have name and url")
		return
	}

	feed, err := v1.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedReq.Name,
		Url:       feedReq.Url,
		UserID:    user.ID,
	})

	if err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, fmt.Sprintf("Unable to create feed %v of %v for user %v", feedReq.Name, feedReq.Url, user.Name))
	}

	helpers.RespondWithJSON(w, http.StatusOK, feed)
}
