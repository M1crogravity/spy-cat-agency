-- name: CreateAgent :one
INSERT INTO agents (
  name, password_hash
) VALUES (
  $1, $2
)
RETURNING id;

-- name: FindAgentByName :one
SELECT *
FROM agents
WHERE name = $1
LIMIT 1;

-- name: FindAgentById :one
SELECT *
FROM agents
WHERE id = $1
LIMIT 1;
