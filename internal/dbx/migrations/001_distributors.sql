-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "distributors_status" AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TABLE "distributors" (
    "id"         UUID               PRIMARY KEY NOT NULL,
    "icon"       VARCHAR(256)       NOT NULL,
    "name"       VARCHAR(256)       NOT NULL,
    "status"     distributors_status NOT NULL DEFAULT 'active',
    "updated_at" TIMESTAMP          NOT NULL,
    "created_at" TIMESTAMP          NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS "distributors" CASCADE;
DROP TYPE IF EXISTS "distributors_status" CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
