-- +goose Up
CREATE TABLE services
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR   NOT NULL,
    description VARCHAR   NOT NULL,
    is_removed  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);

CREATE TYPE service_event_type AS ENUM ('CREATED', 'UPDATED', 'REMOVED');
CREATE TYPE service_event_status AS ENUM ('DEFERRED', 'PROCESSED');

CREATE TABLE service_events
(
    id         BIGSERIAL PRIMARY KEY,
    service_id BIGSERIAL            NOT NULL,
    type       service_event_type   NOT NULL,
    status     service_event_status NOT NULL,
    payload    JSONB                NOT NULL,
    updated_at TIMESTAMP            NOT NULL
);

-- +goose Down
DROP TABLE services;
DROP TABLE service_events;
DROP TYPE service_event_type;
DROP TYPE service_event_status;
