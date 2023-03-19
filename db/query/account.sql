-- name: CreateAccount :one
INSERT INTO accounts (
    "account_number",
    "owner",
    "balance",
    "currency"
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE account_number = $1 LIMIT 1;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE account_number = $1
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE account_number = sqlc.arg(account_number)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_number = $1;

