-- name: CreateEntry :one
INSERT INTO entries (
    "account_number",
    "amount"
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntryFromAccountNumber :many
SELECT * FROM entries
WHERE account_number = $1
LIMIT $2
OFFSET $3;

-- name: GetEntryById :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_number = $1
ORDER BY id
    LIMIT $2
OFFSET $3;