-- name: SelectUser :one
SELECT *
FROM user
WHERE login = ?;

-- name: InsertUser :exec
INSERT INTO user (login, password)
VALUES (?, ?);

-- name: InsertLogin :one
INSERT INTO login (name, login, password, website, notes, user)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, updated_at;

-- name: UpdateLogin :one
UPDATE login
SET name     = ?,
    login    = ?,
    password = ?,
    website  = ?,
    notes    = ?,
    updated_at = unixepoch()
WHERE id = ?
RETURNING updated_at;

-- name: SelectLogin :one
SELECT *
FROM login
WHERE id = ?
  AND user = ?;

-- name: SelectLoginUser :one
SELECT user
FROM login
WHERE id = ?;

-- name: SelectLogins :many
SELECT *
FROM login
WHERE user = ?;

-- name: DeleteLogin :execrows
DELETE
FROM login
WHERE id = ?
  AND user = ?;

-- name: InsertNote :one
INSERT INTO note (name, text, user)
VALUES (?, ?, ?)
RETURNING id, updated_at;

-- name: UpdateNote :one
UPDATE note
SET name = ?,
    text = ?,
    updated_at = unixepoch()
WHERE id = ?
RETURNING updated_at;

-- name: SelectNote :one
SELECT *
FROM note
WHERE id = ?
  AND user = ?;

-- name: SelectNoteUser :one
SELECT user
FROM note
WHERE id = ?;

-- name: SelectNotes :many
SELECT *
FROM note
WHERE user = ?;

-- name: DeleteNote :execrows
DELETE
FROM note
WHERE id = ?
  AND user = ?;

-- name: InsertBinary :one
INSERT INTO binary (name, filename, size, notes, user)
VALUES (?, ?, ?, ?, ?)
RETURNING id, updated_at;

-- name: UpdateBinary :one
UPDATE binary
SET name  = ?,
    notes = ?,
    updated_at = unixepoch()
WHERE id = ?
RETURNING updated_at;

-- name: SelectBinary :one
SELECT *
FROM binary
WHERE id = ?
  AND user = ?;

-- name: SelectBinaryUser :one
SELECT user
FROM binary
WHERE id = ?;

-- name: SelectBinaries :many
SELECT *
FROM binary
WHERE user = ?;

-- name: DeleteBinary :execrows
DELETE
FROM binary
WHERE id = ?
  AND user = ?;

-- name: InsertCard :one
INSERT INTO card (name, number, exp_date, cvv, cardholder, notes, user)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id, updated_at;

-- name: UpdateCard :one
UPDATE card
SET name       = ?,
    number     = ?,
    exp_date   = ?,
    cvv        = ?,
    cardholder = ?,
    notes      = ?,
    updated_at = unixepoch()
WHERE id = ?
RETURNING updated_at;

-- name: SelectCard :one
SELECT *
FROM card
WHERE id = ?
  AND user = ?;

-- name: SelectCardUser :one
SELECT user
FROM card
WHERE id = ?;

-- name: SelectCards :many
SELECT *
FROM card
WHERE user = ?;

-- name: DeleteCard :execrows
DELETE
FROM card
WHERE id = ?
  AND user = ?;
