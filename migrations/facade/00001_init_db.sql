-- +goose Up
CREATE TABLE services
(
    id            BIGINT PRIMARY KEY,
    name          VARCHAR   NOT NULL,
    description   VARCHAR   NOT NULL,
    is_removed    BOOLEAN   NOT NULL DEFAULT FALSE,
    last_event_id BIGINT NOT NULL
);

CREATE INDEX services_is_removed_idx ON services USING btree (is_removed);
CREATE UNIQUE INDEX services_last_event_id_idx ON services USING btree (last_event_id);

-- +goose Down
DROP INDEX services_last_event_id_idx;
DROP INDEX services_is_removed_idx;
DROP TABLE services;
