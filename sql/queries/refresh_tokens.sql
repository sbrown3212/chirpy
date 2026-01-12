-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  token,
  create_at,
  updated_at,
  user_id,
  expires_at
)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2, 
  NOW() + INTERVAL '60 days'
)
RETURNING *;
