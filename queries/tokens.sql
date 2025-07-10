-- name: CreateToken :exec
INSERT INTO tokens (
  hash,
  user_id,
  user_type,
  expiry,
  scope
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: FindTokenByPlaintext :one
SELECT *
FROM tokens
WHERE hash = $1
  AND scope = $2
  AND expiry >= $3;
