-- +goose Up
CREATE TABLE service (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE service;
