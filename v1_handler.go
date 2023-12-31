package main

import (
	"net/http"

	v1handlers "github.com/AtinAgnihotri/glog-aggregator/v1_handlers"
	"github.com/go-chi/chi/v5"
)

const (
	readiness    = "/readiness"
	err          = "/err"
	users        = "/users"
	feeds        = "/feeds"
	feed_follows = "/feed_follows"
	posts        = "/posts"
)

func V1Handler(serverConf *ServerConf) http.Handler {
	r := chi.NewRouter()

	v1Handlers := v1handlers.V1Handlers{
		DB: serverConf.DB,
	}

	go v1handlers.RssWorker(&v1Handlers)

	// health check endpoints
	r.Get(readiness, v1handlers.Readiness)
	r.Get(err, v1handlers.Err)

	// /users endpoint
	r.Post(users, v1Handlers.CreateUser)
	r.Get(users, v1Handlers.MiddlewareAuth(v1Handlers.GetUser))

	// /feeds endpoint
	r.Post(feeds, v1Handlers.MiddlewareAuth(v1Handlers.CreateFeed))
	r.Get(feeds, v1Handlers.GetAllFeeds)

	// /feed_follows endpoint
	r.Post(feed_follows, v1Handlers.MiddlewareAuth(v1Handlers.FollowFeed))
	r.Delete(feed_follows+"/{feedFollowID}", v1Handlers.MiddlewareAuth(v1Handlers.RemoveFeedFollowing))
	r.Get(feed_follows, v1Handlers.MiddlewareAuth(v1Handlers.GetFeedFollowingsByUser))

	// /posts endpoint
	r.Get(posts, v1Handlers.MiddlewareAuth(v1Handlers.GetPosts))
	r.Get(posts+"/{limit}", v1Handlers.MiddlewareAuth(v1Handlers.GetPosts))

	return r
}
