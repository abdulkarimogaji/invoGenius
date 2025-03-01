-- name: PingDB :one
SELECT NOW();

-- name: CreateUser :execresult
INSERT INTO user (first_name, last_name, role, email, password, created_at, updated_at) VALUES (?,?,?,?,?,?,?);

-- name: GetUserByEmail :one
SELECT * FROM user WHERE email = ?;