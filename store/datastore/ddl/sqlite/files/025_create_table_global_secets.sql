-- name: create-table-global-secrets

CREATE TABLE IF NOT EXISTS global_secrets (
 secret_id          INTEGER PRIMARY KEY AUTOINCREMENT
,secret_name        TEXT
,secret_value       TEXT
,secret_images      TEXT
,secret_events      TEXT
,secret_skip_verify BOOLEAN
,secret_conceal     BOOLEAN
,UNIQUE(secret_name)
);

-- name: create-index-global-secrets-name

CREATE INDEX IF NOT EXISTS ix_globalsecrets_name ON global_secrets (secret_name);
