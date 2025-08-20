-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE suspended_distributor_status AS ENUM (
    'active',   -- distributor is currently suspended
    'cancelled' -- suspension has been cancelled
);

CREATE TABLE suspended_distributors (
    id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    distributor_id uuid NOT NULL REFERENCES distributors(id) ON DELETE CASCADE,
    initiator_id   uuid NOT NULL,
    reason         varchar(8192) NOT NULL,
    status         suspended_distributor_status NOT NULL DEFAULT 'active',
    suspended_at   timestamp NOT NULL DEFAULT now(),
    cancelled_at   timestamp,
    created_at     timestamp NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX suspended_distributors_one_active_per_dist
    ON suspended_distributors(distributor_id)
    WHERE status = 'active';

-- +migrate Down
DROP TABLE IF EXISTS "suspended_distributors" CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;
DROP INDEX IF EXISTS suspended_distributors_one_active_per_dist;