-- name: GetUser :one
SELECT
  *
FROM
  user
WHERE
  id = ?
LIMIT
  1;