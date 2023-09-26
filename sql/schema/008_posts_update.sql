-- +goose Up

ALTER TABLE posts 
  ADD CONSTRAINT feed_id
  FOREIGN KEY (feed_id)
  REFERENCES feeds(id)
  ON DELETE CASCADE;

-- +goose Down

ALTER TABLE posts 
  DROP CONSTRAINT feed_id;
