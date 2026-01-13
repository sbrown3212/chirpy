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

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :execrows
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;
