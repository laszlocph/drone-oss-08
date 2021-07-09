-- name: create-table-global-secrets

CREATE TABLE IF NOT EXISTS global_secrets (
 secret_id          INTEGER PRIMARY KEY AUTO_INCREMENT
,secret_name        VARCHAR(250)
,secret_value       MEDIUMBLOB
,secret_images      VARCHAR(2000)
,secret_events      VARCHAR(2000)
,secret_skip_verify BOOLEAN
,secret_conceal     BOOLEAN

,UNIQUE(secret_name)
);

-- name: create-index-global-secrets-name

CREATE INDEX ix_globalsecrets_name  ON global_secrets  (secret_name);
