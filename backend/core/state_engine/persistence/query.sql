-- name: FindStateByEmail :one
SELECT email, state FROM conversation_state WHERE email = $1;

-- name: UpdateState :exec
INSERT INTO conversation_state (email, state)
VALUES ($1, $2)
ON CONFLICT (email) DO UPDATE SET
    state = EXCLUDED.state;