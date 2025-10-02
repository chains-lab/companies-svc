-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE distributors_status AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TABLE distributors (
    id         UUID                PRIMARY KEY NOT NULL,
    icon       VARCHAR(256)        NOT NULL,
    name       VARCHAR(256)        NOT NULL,
    status     distributors_status NOT NULL DEFAULT 'active',
    updated_at TIMESTAMP           NOT NULL,
    created_at TIMESTAMP           NOT NULL
);

CREATE TYPE blocked_distributor_status AS ENUM (
    'active',   -- distributor is currently blocked
    'cancelled' -- suspension has been cancelled
);

CREATE TABLE distributor_blockages (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    distributor_id UUID NOT NULL REFERENCES distributors("id") ON DELETE CASCADE,
    initiator_id   UUID NOT NULL,
    reason         varchar(8192) NOT NULL,
    status         blocked_distributor_status NOT NULL DEFAULT 'active',
    blocked_at     timestamp NOT NULL DEFAULT now(),
    cancelled_at   timestamp
);

CREATE UNIQUE INDEX blocked_distributors_one_active_per_dist
    ON distributor_blockages(distributor_id)
    WHERE status = 'active';

-- +migrate Down
DROP TABLE IF EXISTS distributors CASCADE;
DROP TABLE IF EXISTS distributor_blockages CASCADE;

DROP TYPE IF EXISTS distributors_status CASCADE;
DROP TYPE IF EXISTS blocked_distributor_status CASCADE;

DROP INDEX IF EXISTS blocked_distributors_one_active_per_dist;

DROP EXTENSION IF EXISTS "uuid-ossp";
