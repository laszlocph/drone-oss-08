-- name: global-secret-find

SELECT
 secret_id
,secret_name
,secret_value
,secret_images
,secret_events
,secret_conceal
,secret_skip_verify
FROM global_secrets

-- name: global-secret-find-name

SELECT
secret_id
,secret_name
,secret_value
,secret_images
,secret_events
,secret_conceal
,secret_skip_verify
FROM global_secrets
WHERE secret_name = $1

-- name: global-secret-delete

DELETE FROM global_secrets WHERE secret_id = $1
