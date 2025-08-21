-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE blocked_distributor_status AS ENUM (
    'active',   -- distributor is currently blocked
    'cancelled' -- suspension has been cancelled
);

CREATE TABLE blocked_distributors (
    id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    distributor_id uuid NOT NULL REFERENCES distributors(id) ON DELETE CASCADE,
    initiator_id   uuid NOT NULL,
    reason         varchar(8192) NOT NULL,
    status         locked_distributor_status NOT NULL DEFAULT 'active',
    blocked_at     timestamp NOT NULL DEFAULT now(),
    cancelled_at   timestamp,
    created_at     timestamp NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX blocked_distributors_one_active_per_dist
    ON blocked_distributors(distributor_id)
    WHERE status = 'active';

-- +migrate Down
DROP TABLE IF EXISTS "blocked_distributors" CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;
DROP INDEX IF EXISTS blocked_distributors_one_active_per_dist;