package v1handlers

import (
	"fmt"
	"net/http"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (v1 *V1Handlers) FollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowReq := FeedFollowRequest{}
	err := helpers.DecodeBodyToJson[FeedFollowRequest](r, &feedFollowReq)

	if err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, "Invalid Request. Body needs to have feed_id")
		return
	}

	feedFollow, err := v1.followFeed(w, r, user, feedFollowReq.FeedId)

	if err != nil {
		return
	} // Error response already handled

	helpers.RespondWithJSON(w, http.StatusOK, FeedFollowResponse{
		Id:        feedFollow.ID,
		UserId:    user.ID,
		FeedId:    feedFollow.FeedID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	})
}

func (v1 *V1Handlers) RemoveFeedFollowing(w http.ResponseWriter, r *http.Request, user database.User) {
	param := chi.URLParam(r, "feedFollowID")
	feedFollowId, err := uuid.Parse(param)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request. Need a valid feed following id")
		return
	}

	err = v1.DB.RemoveFollowFeed(r.Context(), feedFollowId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to remove feed following from %v", feedFollowId))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (v1 *V1Handlers) GetFeedFollowingsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowings, err := v1.DB.GetFollowFeedByUserId(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Unable to find feed following for user %v", user.Name))
	}

	var feedFollows []FeedFollowResponse
	for _, feedFollowing := range feedFollowings {
		feedFollows = append(feedFollows, FeedFollowResponse{
			Id:        feedFollowing.ID,
			FeedId:    feedFollowing.FeedID,
			UserId:    feedFollowing.UserID,
			CreatedAt: feedFollowing.CreatedAt,
			UpdatedAt: feedFollowing.UpdatedAt,
		})
	}

	helpers.RespondWithJSON(w, http.StatusOK, feedFollows)
}
