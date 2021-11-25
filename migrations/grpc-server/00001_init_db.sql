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

CREATE INDEX services_is_removed_idx ON services USING btree (is_removed);

CREATE TYPE service_event_type AS ENUM ('CREATED', 'UPDATED', 'REMOVED');
CREATE TYPE service_event_subtype AS ENUM ('NONE', 'NAME', 'DESCRIPTION');
CREATE TYPE service_event_status AS ENUM ('DEFERRED', 'PROCESSED');

CREATE TABLE service_events
(
    id         BIGSERIAL PRIMARY KEY,
    service_id BIGINT                NOT NULL,
    type       service_event_type    NOT NULL,
    subtype    service_event_subtype NOT NULL DEFAULT 'NONE',
    status     service_event_status  NOT NULL,
    payload    JSONB                 NOT NULL,
    updated_at TIMESTAMP             NOT NULL
);

CREATE INDEX service_events_service_id_idx ON service_events USING btree (service_id);
CREATE INDEX service_events_type_idx ON service_events USING btree (type);
CREATE INDEX service_events_subtype_idx ON service_events USING btree (subtype);
CREATE INDEX service_events_status_idx ON service_events USING btree (status);

-- +goose Down
DROP INDEX services_is_removed_idx;
DROP INDEX service_events_status_idx;
DROP INDEX service_events_type_idx;
DROP INDEX service_events_service_id_idx;
DROP TABLE services;
DROP TABLE service_events;
DROP TYPE service_event_type;
DROP TYPE service_event_subtype;
DROP TYPE service_event_status;
