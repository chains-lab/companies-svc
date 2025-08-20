-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "distributor_status" AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE "distributors" (
    "id"         UUID               PRIMARY KEY NOT NULL,
    "icon"       VARCHAR(256)       NOT NULL,
    "name"       VARCHAR(256)       NOT NULL,
    "status"     distributor_status NOT NULL DEFAULT 'active',
    "updated_at" TIMESTAMP          NOT NULL,
    "created_at" TIMESTAMP          NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS "distributors" CASCADE;
DROP TYPE IF EXISTS "distributor_status" CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
