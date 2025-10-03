-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE companies_status AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TABLE companies (
    id         UUID                PRIMARY KEY NOT NULL,
    icon       VARCHAR(256)        NOT NULL,
    name       VARCHAR(256)        NOT NULL,
    status     companies_status NOT NULL DEFAULT 'active',
    updated_at TIMESTAMP           NOT NULL,
    created_at TIMESTAMP           NOT NULL
);

CREATE TYPE blocked_company_status AS ENUM (
    'active',   -- company is currently blocked
    'cancelled' -- suspension has been cancelled
);

CREATE TABLE company_blocks (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies("id") ON DELETE CASCADE,
    initiator_id   UUID NOT NULL,
    reason         varchar(8192) NOT NULL,
    status         blocked_company_status NOT NULL DEFAULT 'active',
    blocked_at     timestamp NOT NULL DEFAULT now(),
    cancelled_at   timestamp
);

CREATE UNIQUE INDEX blocked_companies_one_active_per_dist
    ON company_blocks(company_id)
    WHERE status = 'active';

-- +migrate Down
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS company_blocks CASCADE;

DROP TYPE IF EXISTS companies_status CASCADE;
DROP TYPE IF EXISTS blocked_company_status CASCADE;

DROP INDEX IF EXISTS blocked_companies_one_active_per_dist;

DROP EXTENSION IF EXISTS "uuid-ossp";
