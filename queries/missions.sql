-- name: CreateMission :one
INSERT INTO missions (
  state
) VALUES (
  $1
) RETURNING id;

-- name: CreateTargets :copyfrom
INSERT INTO targets (
  id,
  mission_id,
  name,
  country,
  notes,
  state
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: FindMissionById :many
SELECT missions.id as mission_id,
  missions.state as mission_state,
  missions.spy_cat_id,
  targets.id as target_id,
  targets.name,
  targets.country,
  targets.notes,
  targets.state as target_state
FROM missions
INNER JOIN targets ON missions.id = targets.mission_id
WHERE missions.id = $1;

-- name: DeleteMission :exec
DELETE
FROM missions
WHERE id = $1;

-- name: UpdateMission :exec
UPDATE missions
SET state = $2,
  spy_cat_id = $3
WHERE id = $1;

-- name: UpdateTarget :exec
UPDATE targets
SET notes = $3,
  state = $4
WHERE id = $1
  AND mission_id = $2;

-- name: DeleteTarget :exec
DELETE
FROM targets
WHERE id = $1
  AND mission_id = $2;

-- name: FindActiveMission :many
SELECT missions.id as mission_id,
  missions.state as mission_state,
  missions.spy_cat_id,
  targets.id as target_id,
  targets.name,
  targets.country,
  targets.notes,
  targets.state as target_state
FROM missions
INNER JOIN targets ON missions.id = targets.mission_id
WHERE missions.spy_cat_id = $1
  AND missions.state = 'in_progress';

-- name: FindAllMissions :many
SELECT missions.id as mission_id,
  missions.state as mission_state,
  missions.spy_cat_id,
  targets.id as target_id,
  targets.name,
  targets.country,
  targets.notes,
  targets.state as target_state
FROM missions
INNER JOIN targets ON missions.id = targets.mission_id
ORDER BY missions.id ASC;
