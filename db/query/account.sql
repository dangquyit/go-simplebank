-- name: CreateAccount :one
INSERT INTO accounts (
    "account_number",
    "owner",
    "user_name",
    "email",
    "password",
    "phone_number",
    "balance",
    "currency"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE account_number = $1 LIMIT 1;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE account_number = $1
RETURNING account_number, balance, created_at, updated_at;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_number = $1;

