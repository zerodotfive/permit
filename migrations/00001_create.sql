-- +goose Up
CREATE TABLE permits (
                         id                              uuid PRIMARY KEY,
                         parent_id                       uuid REFERENCES permits(id) ON DELETE CASCADE,
                         comment                         text NOT NULL CHECK (comment <> ''),
                         embed_site                      text NOT NULL CHECK (embed_site <> ''),
                         target_site                     text NOT NULL DEFAULT '',
                         target_location                 text NOT NULL DEFAULT '',
                         token_sha                       bytea UNIQUE NOT NULL,
                         token_use_not_after             timestamptz NOT NULL,
                         session_lifetime_minutes        integer NOT NULL CHECK (session_lifetime_minutes BETWEEN 0 AND 1440),
                         session_token_sha               bytea UNIQUE,
                         session_token_use_not_after     timestamptz,
                         CHECK (parent_id IS NULL OR parent_id <> id),
                         CHECK (parent_id IS NOT NULL OR session_lifetime_minutes >= 1)
);

CREATE INDEX permits_parent_id_idx ON permits (parent_id);

-- +goose Down
DROP TABLE permits;
