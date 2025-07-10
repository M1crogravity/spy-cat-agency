-- name: CreateSpyCat :one
INSERT INTO spy_cats (
  name,
  password_hash,
  years_of_experience,
  breed,
  salary
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: FindSpyCatById :one
SELECT *
FROM spy_cats
WHERE id = $1
LIMIT 1;

-- name: DeleteSpyCatById :exec
DELETE
FROM spy_cats
WHERE id = $1;

-- name: UpdateSpyCat :exec
UPDATE spy_cats
SET salary = $2
WHERE id = $1;

-- name: ListSpyCats :many
SELECT *
FROM spy_cats
ORDER BY id ASC;

-- name: FindSpyCatByName :one
SELECT *
FROM spy_cats
WHERE name = $1
LIMIT 1;
