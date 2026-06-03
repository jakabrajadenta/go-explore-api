-- =============================================================
-- DDL: go_explore database
-- Schema : user_management
-- Run    : Execute via DBeaver, psql, or any SQL client
-- =============================================================

-- 1. Create schema
CREATE SCHEMA IF NOT EXISTS user_management;

-- =============================================================
-- 2. Tables
-- =============================================================

CREATE TABLE IF NOT EXISTS user_management.users (
    id          BIGSERIAL       NOT NULL,
    username    VARCHAR(50)     NOT NULL,
    email       VARCHAR(255)    NOT NULL,
    full_name   VARCHAR(255)    NOT NULL,
    phone       VARCHAR(20)     NOT NULL DEFAULT '',
    is_active   BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_users            PRIMARY KEY (id),
    CONSTRAINT uq_users_username   UNIQUE (username),
    CONSTRAINT uq_users_email      UNIQUE (email)
);

COMMENT ON TABLE  user_management.users            IS 'Registered application users';
COMMENT ON COLUMN user_management.users.username   IS 'Unique login identifier, max 50 chars';
COMMENT ON COLUMN user_management.users.email      IS 'Unique email address';
COMMENT ON COLUMN user_management.users.full_name  IS 'Display / full legal name';
COMMENT ON COLUMN user_management.users.phone      IS 'Optional phone number, empty string when not provided';
COMMENT ON COLUMN user_management.users.is_active  IS 'Soft-enable/disable flag; TRUE = active';

-- =============================================================
-- 3. Indexes
-- =============================================================

CREATE INDEX IF NOT EXISTS idx_users_email     ON user_management.users (email);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON user_management.users (is_active);

-- =============================================================
-- 4. Auto-update updated_at trigger
-- =============================================================

CREATE OR REPLACE FUNCTION user_management.fn_set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_users_updated_at ON user_management.users;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON user_management.users
    FOR EACH ROW
    EXECUTE FUNCTION user_management.fn_set_updated_at();
