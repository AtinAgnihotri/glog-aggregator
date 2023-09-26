package v1handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/helpers"
	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	DEFAULT_LIMIT = 10
)

type PostResponse struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func databasePostToPostResponse(post database.Post) PostResponse {
	var description *string = nil
	var publishedAt *time.Time = nil

	if post.Description.Valid {
		description = &post.Description.String
	}

	if post.PublishedAt.Valid {
		publishedAt = &post.PublishedAt.Time
	}

	return PostResponse{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: description,
		PublishedAt: publishedAt,
		FeedID:      post.FeedID,
	}
}

func databasePostsToPostReponses(posts []database.Post) []PostResponse {
	var resp []PostResponse
	for _, post := range posts {
		resp = append(resp, databasePostToPostResponse(post))
	}
	return resp
}

func (v1 *V1Handlers) GetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	param := chi.URLParam(r, "limit")
	limit, err := strconv.Atoi(param)
	if err != nil {
		limit = DEFAULT_LIMIT
	}
	posts, err := v1.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't find posts for user %v", user.Name))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, databasePostsToPostReponses(posts))

}
