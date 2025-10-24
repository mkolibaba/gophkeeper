-- name: SelectUser :one
SELECT *
FROM user
WHERE login = ?;

-- name: InsertUser :exec
INSERT INTO user (login, password)
VALUES (?, ?);

-- name: InsertLogin :execlastid
INSERT INTO login (name, login, password, website, notes, user)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateLogin :execrows
UPDATE login
SET name     = ?,
    login    = ?,
    password = ?,
    website  = ?,
    notes    = ?
WHERE id = ?;

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

-- name: InsertNote :execlastid
INSERT INTO note (name, text, user)
VALUES (?, ?, ?);

-- name: UpdateNote :execrows
UPDATE note
SET name = ?,
    text = ?
WHERE id = ?;

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
RETURNING id;

-- name: UpdateBinary :execrows
UPDATE binary
SET name  = ?,
    notes = ?
WHERE id = ?;

-- name: SelectBinary :one
SELECT *
FROM binary
WHERE id = ?;

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
WHERE id = ?;

-- name: InsertCard :execlastid
INSERT INTO card (name, number, exp_date, cvv, cardholder, notes, user)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateCard :execrows
UPDATE card
SET name       = ?,
    number     = ?,
    exp_date   = ?,
    cvv        = ?,
    cardholder = ?,
    notes      = ?
WHERE id = ?;

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
