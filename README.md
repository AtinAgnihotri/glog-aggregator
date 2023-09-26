# glog-aggregator

A simple blog aggregator server written in go. Users can sign up, register to follow a feed, and fetch the latest posts from that feed. A long-running service worker, running in it's own goroutine, fetches the latest posts for the feeds fetched lasts.

## Setup

1. Clone this repo
2. Setup a .env file with the following:

```shell
PORT="<port>" # I use 8080
DB_URL="postgres://<username>:<pwd>@<posgres_server_host>:<postgress_port>/<db_name>?sslmode=disable"
```

3. Hit the following to start the server:

```shell
sh runServer.sh
```

4. Use your favorite HTTP client to play around with the endpoints.

## Endpoints

Those marked with `A*` are authorized endpoint. They require Authorization header to be present in the request in the following shape:

```
ApiKey <api_key>
```

More on how ApiKey is retrived below.

### `Health-checks endpoints`

1. `/v1/readiness`: Health-check endpoint
2. `/err`: Dummy endpoint. Always returns HTTP 500, Internal Server Error

### `/v1/users`

1. [`POST`]: Requires name of the user in body.

Sample Payload:

```json
{
  "name": "Jatin"
}
```

Sample Response:

```
{
  "ID": "29d83058-b3dd-40e0-80fd-edb0b6efa783",
  "CreatedAt": "2023-09-26T22:54:32.911213Z",
  "UpdatedAt": "2023-09-26T22:54:32.911213Z",
  "Name": "Jatin",
  "ApiKey": "a5fac996c986e268c72ceb7b343b2c120a378041de711c3bca1870a84b3f4cfc"
}
```

2. [`GET`] `A*`: Retrives User info. Requires auth header.

Sample Response:

```json
{
  "ID": "ea474a08-69d2-47c7-93f7-90758231f5f3",
  "CreatedAt": "2023-09-23T03:38:30.81462Z",
  "UpdatedAt": "2023-09-23T03:38:30.81462Z",
  "Name": "Atin",
  "ApiKey": "5487aa03f979c556ce1c1d689bc622c65ab267ddaf680ae6284ac108a1de97ab"
}
```

### `/v1/feeds`

[`POST`] `A*`: Sets up a new feed. Requires auth header. Automatically makes the requesting user a follower of the feed.

Sample Payload:

```json
{
  "name": "Atin's RSS Feed",
  "url": "https://example.atin.com/index.xml"
}
```

Sample Response:

```json
{
  "feed": {
    "id": "ba3b2ca0-66b6-4a79-b7ee-ce50b4f00562",
    "created_at": "2023-09-26T22:59:44.290519Z",
    "updated_at": "2023-09-26T22:59:44.290519Z",
    "name": "WagsLane's RSS Feed",
    "url": "https://example.wagslane.com/index.xml",
    "user_id": "ea474a08-69d2-47c7-93f7-90758231f5f3"
  },
  "feed_follow": {
    "id": "b09491aa-7eb5-4e07-b272-41e96ff29082",
    "feed_id": "ba3b2ca0-66b6-4a79-b7ee-ce50b4f00562",
    "user_id": "ea474a08-69d2-47c7-93f7-90758231f5f3",
    "created_at": "2023-09-26T22:59:44.291047Z",
    "updated_at": "2023-09-26T22:59:44.291047Z"
  }
}
```

### `/v1/feed_follows`

1. [`POST`] `A*`: Makes the requesting user a follower of the feed. Requires auth header.

Sample Payload:

```json
{
  "feed_id": "1533df7a-78d9-46fd-abf4-9004c1230f6e"
}
```

Sample Response:

```json
{
  "id": "c6196423-9f2a-453f-ab8a-c99aa7fe027b",
  "feed_id": "1533df7a-78d9-46fd-abf4-9004c1230f6e",
  "user_id": "ea474a08-69d2-47c7-93f7-90758231f5f3",
  "created_at": "2023-09-26T21:25:32.075131Z",
  "updated_at": "2023-09-26T21:25:32.075131Z"
}
```

2. [`DEL`] `A*`: Removes the requesting user from following the specified feed. Requires auth header and takes the feed id as a param.

Sample request URL:

`/v1/feed_follows/2a9cd973-96c5-401a-87fe-5ed6a7e33c01`

3. [`GET`] `A*`: Returns all the feeds followed by the user. Requires auth header.

Sample Response

```json
[
  {
    "id": "d8af7f0f-31c2-45d9-870b-caa663ee612c",
    "feed_id": "ed11a852-1446-44ab-9c1a-c024739988fc",
    "user_id": "ea474a08-69d2-47c7-93f7-90758231f5f3",
    "created_at": "2023-09-24T23:13:15.594429Z",
    "updated_at": "2023-09-24T23:13:15.594429Z"
  },
  {
    "id": "cad6a960-9186-41f5-a881-8cd050832d56",
    "feed_id": "334cb6a7-383d-4234-b20d-7d9b5ec7fb0d",
    "user_id": "ea474a08-69d2-47c7-93f7-90758231f5f3",
    "created_at": "2023-09-24T23:13:36.363979Z",
    "updated_at": "2023-09-24T23:13:36.363979Z"
  }
]
```

### `/v1/posts`

1. [`GET`] `A*`: Gets the latests posts from the feeds followed by the requesting user. Requires auth header and takes an optional limit url parameter, that can limit how many posts are fetched.

Sample url:

`/v1/posts/2`: Will fetch 2 latest posts

If limit is not defined, or incorrectly sent, fetches the 10 latest posts
