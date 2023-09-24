-- +goose Up

ALTER TABLE feed_follows ADD COLUMN 
created_at TIMESTAMP NOT NULL;
ALTER TABLE feed_follows ADD COLUMN 
updated_at TIMESTAMP NOT NULL;

-- +goose Down

ALTER TABLE feed_follows DROP COLUMN created_at;
ALTER TABLE feed_follows DROP COLUMN updated_at;