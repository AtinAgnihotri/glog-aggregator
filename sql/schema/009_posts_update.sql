-- +goose Up

ALTER TABLE posts
ADD CONSTRAINT url UNIQUE(url);

-- +goose Down

ALTER TABLE posts
DROP CONSTRAINT url;