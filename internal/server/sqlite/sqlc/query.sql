-- name: SaveLogin :exec
INSERT INTO login (name, login, password, metadata, user)
VALUES (?, ?, ?, ?, ?);

-- name: GetAllLogins :many
SELECT *
FROM login
WHERE user = ?;

-- name: RemoveLogin :execrows
DELETE
FROM login
WHERE name = ?;

-- name: SaveNote :exec
INSERT INTO note (name, text, metadata, user)
VALUES (?, ?, ?, ?);

-- name: GetAllNotes :many
SELECT *
FROM note
WHERE user = ?;

-- name: RemoveNote :execrows
DELETE
FROM note
WHERE name = ?;

-- name: SaveBinary :exec
INSERT INTO binary (name, filename, metadata, user)
VALUES (?, ?, ?, ?);

-- name: GetBinary :one
SELECT *
FROM binary
WHERE name = ?
  AND user = ?;

-- name: GetAllBinaries :many
SELECT *
FROM binary
WHERE user = ?;

-- name: RemoveBinary :execrows
DELETE
FROM binary
WHERE name = ?;

-- name: SaveCard :exec
INSERT INTO card (name, number, exp_date, cvv, cardholder, metadata, user)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetAllCards :many
SELECT *
FROM card
WHERE user = ?;

-- name: RemoveCard :execrows
DELETE
FROM card
WHERE name = ?;
