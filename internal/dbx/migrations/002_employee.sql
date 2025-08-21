-- +migrate Up

CREATE TYPE "employee_roles" AS ENUM (
    'owner',
    'admin',
    'moderator'
);

CREATE TABLE "employees" (
    "user_id"         UUID           PRIMARY KEY NOT NULL,
    "distributor_id"  UUID           NOT NULL,
    "role"            employee_roles NOT NULL,
    "updated_at"      TIMESTAMP      NOT NULL,
    "created_at"      TIMESTAMP      NOT NULL
);

CREATE Type "invite_status" AS ENUM (
    'sent',
    'recalled',
    'accepted',
    'rejected'
);

CREATE TABLE "employee_invites" (
    "id"              UUID           PRIMARY KEY NOT NULL,
    "distributor_id"  UUID           NOT NULL REFERENCES "distributors" ("id") ON DELETE CASCADE,
    "user_id"         UUID           NOT NULL,
    "invited_by"      UUID           NOT NULL,
    "role"            employee_roles NOT NULL,
    "status"          invite_status  NOT NULL DEFAULT 'sent',
    "answered_at"     TIMESTAMP      NULL,
    "expires_at"      TIMESTAMP      NOT NULL,
    "created_at"      TIMESTAMP      NOT NULL

    CHECK (
        (status = 'sent'  AND answered_at IS NULL) OR
        (status IN ('accepted','rejected') AND answered_at IS NOT NULL)
    ),
    CHECK (expires_at > created_at)
);

-- +migrate Down
DROP TABLE IF EXISTS "employee_invites" CASCADE;
DROP TYPE IF EXISTS "invite_status" CASCADE;

DROP TABLE IF EXISTS "employees" CASCADE;
DROP TYPE IF EXISTS "employee_roles" CASCADE;