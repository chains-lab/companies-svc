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
    'accepted',
    'rejected'
);

CREATE TABLE employee_invites (
    "id"             UUID           PRIMARY KEY,
    "status"         invite_status  NOT NULL DEFAULT 'sent',
    "role"           employee_roles NOT NULL,
    "distributor_id" UUID           NOT NULL REFERENCES distributors("id") ON DELETE CASCADE,
    "user_id"        UUID,
    "answered_at"    TIMESTAMP,
    "expires_at"     TIMESTAMP      NOT NULL,
    "created_at"     TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CONSTRAINT invite_status_answered_ck CHECK (
        ("status" = 'sent' AND "answered_at" IS NULL)
        OR ("status" IN ('accepted','rejected') AND "answered_at" IS NOT NULL)
    )
);

-- +migrate Down
DROP TABLE IF EXISTS "employee_invites" CASCADE;
DROP TYPE IF EXISTS "invite_status" CASCADE;

DROP TABLE IF EXISTS "employees" CASCADE;
DROP TYPE IF EXISTS "employee_roles" CASCADE;