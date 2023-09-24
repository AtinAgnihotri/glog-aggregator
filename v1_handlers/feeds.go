package v1handlers

import (
	"fmt"
	"net/http"

	// "strconv"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/google/uuid"
)

type FeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type FeedFollowRequest struct {
	FeedId uuid.UUID `json:"feed_id"`
}

type FeedResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}

type FeedFollowResponse struct {
	Id        uuid.UUID `json:"id"`
	FeedId    uuid.UUID `json:"feed_id"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeedCreationResponse struct {
	Feed       FeedResponse       `json:"feed"`
	FeedFollow FeedFollowResponse `json:"feed_follow"`
}

func (v1 *V1Handlers) followFeed(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	feedId uuid.UUID) (database.FeedFollow, error) {
	feedFollow, err := v1.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feedId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to follow feed %v", feedId))
		return database.FeedFollow{}, err
	}

	return feedFollow, nil
}

func (v1 *V1Handlers) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := v1.DB.GetFeeds(r.Context())

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to fetch feeds")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, feeds)
}

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

	feedFollow, err := v1.followFeed(w, r, user, feed.ID)

	if err != nil {
		return
	} // Error response already handled

	helpers.RespondWithJSON(w, http.StatusOK, FeedCreationResponse{
		Feed: FeedResponse{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		},
		FeedFollow: FeedFollowResponse{
			Id:        feedFollow.ID,
			FeedId:    feedFollow.FeedID,
			UserId:    feedFollow.UserID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
		},
	})
}
