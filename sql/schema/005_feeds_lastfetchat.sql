-- +goose Up
ALTER TABLE feeds ADD COLUMN lastFetchAt TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN lastFetchAt;