-- name: create-table-global-secrets

CREATE TABLE IF NOT EXISTS global_secrets (
 secret_id          SERIAL PRIMARY KEY
,secret_name        VARCHAR(250)
,secret_value       BYTEA
,secret_images      VARCHAR(2000)
,secret_events      VARCHAR(2000)
,secret_skip_verify BOOLEAN
,secret_conceal     BOOLEAN

,UNIQUE(secret_name)
);

-- name: create-index-global-secrets-name

CREATE INDEX IF NOT EXISTS ix_globalsecrets_name ON global_secrets  (secret_name);
