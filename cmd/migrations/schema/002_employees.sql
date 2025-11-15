-- +migrate Up
CREATE TYPE employee_roles AS ENUM (
    'owner',
    'admin',
    'moderator'
);

CREATE TABLE employees (
    user_id    UUID           PRIMARY KEY NOT NULL,
    company_id UUID           NOT NULL REFERENCES companies("id") ON DELETE CASCADE,
    role       employee_roles NOT NULL,
    updated_at TIMESTAMP      NOT NULL  DEFAULT (now() AT TIME ZONE 'UTC'),
    created_at TIMESTAMP      NOT NULL  DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE Type invite_status AS ENUM (
    'sent',
    'accepted',
    'declined'
);

CREATE TABLE invites (
    id           UUID           PRIMARY KEY,
    company_id   UUID           NOT NULL REFERENCES companies("id") ON DELETE CASCADE,
    user_id      UUID           NOT NULL,
    initiator_id UUID           NOT NULL,
    status       invite_status  NOT NULL DEFAULT 'sent',
    role         employee_roles NOT NULL,

    expires_at  TIMESTAMP      NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
);

-- +migrate Down
DROP TABLE IF EXISTS invites CASCADE;
DROP TYPE IF EXISTS invite_status CASCADE;

DROP TABLE IF EXISTS employees CASCADE;
DROP TYPE IF EXISTS employee_roles CASCADE;