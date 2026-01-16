# Chirpy API

Chirpy is a basic Twitter-style backend API built in Go with PostgreSQL. It’s
designed as a learning project that demonstrates how to structure a RESTful
service with user accounts, JWT-based access authentication, and CRUD operations
for short posts (“chirps”). The codebase uses `sqlc` for type-safe database
queries and `goose` for migrations, making it a clear, approachable reference
for how to wire Go and Postgres together in a production-style layout.

## Learning Project Disclaimer

This project was built as part of a guided learning course from
[Boot.dev](https://boot.dev) and is not intended for production use. See
[Tradeoffs / Limitations](#tradeoffs--limitations) for details on what has been
simplified.

## Features

- Create and update user accounts
- Create new chirps
- List all chirps or get a chirp by id
- Uses JWT access tokens and long-lived, opaque, database-backed refresh tokens
for authentication

## Tech Stack

- Backend: Go
- Database: PostgreSQL
- Migrations: `goose`
- DB layer: `sqlc` (type-safe queries)
- Auth: JWT access tokens and opaque refresh tokens
- Password hashing: `Argon2id`
- Configuration: `godotenv` (loading environment variables)

## Authentication Overview

Chirpy uses JSON Web Tokens (JWTs) for authenticated requests and a separate
refresh token to obtain new access tokens without re-entering credentials.

- **Access tokens (JWTs)**
  - Issued on login.
  - Included in the `Authorization: Bearer <access_token>` header for protected
  endpoints.  
  - Contain user identity and an expiration time.
  - Are not stored server-side; they are verified using a shared secret.

- **Refresh tokens (opaque strings)**
  - Issued on login along with the access token.
  - Long-lived, 256-bit random strings stored in PostgreSQL.
  - Used with `POST /api/refresh` to obtain a new access token.
  - Can be revoked with `POST /api/revoke`, which invalidates the stored token.

This design keeps access tokens short-lived and stateless while allowing
sessions to continue using refresh tokens, which can be revoked by the server as
needed.

## Tradeoffs / Limitations

Chirpy’s authentication system is intentionally simplified for learning purposes:

- Refresh tokens are long-lived and are **not rotated** on use (new refresh
token is only issued at login).
- Refresh tokens are stored server-side in PostgreSQL and can be revoked, but
there is no detection of token reuse or theft.
- The API is designed for local development and educational use only and has not
been hardened for production (rate limiting, auditing, monitoring, etc. are out
of scope).

These tradeoffs keep the codebase small and focused on core concepts like
RESTful routing, JWT access tokens, and basic token revocation.

## Setup

### Requirements

- Go 1.23+
- PostgreSQL 15+
- `psql` command-line client
- `goose` (migrations)

### Create Database

Create a database for Chirpy:

```sh
createdb chirpy
```

### Clone Repository

```sh
git clone https://github.com/sbrown3212/chirpy.git
```

### Environment Variables

There should be a file named `.env.example` at the root of the project. This
file shows the format the Chirpy server expects for the `.env` file, along with
explanations of what values should be used.

Rename this `.env.example` file to `.env`, or copy it to a file named `.env`.
Then, update the variables as needed.

> Be sure not to commit your `.env` file to version control.

### Run Database Migrations

Load environment variables for CLI tools (like `goose`):

```sh
source .env
```

Run migrations (uses `$DB_URL` from `.env`):

```sh
goose postgres "$DB_URL" up
```

### Run Locally

From the project root, run:

```sh
go run .
```

## Base URL

Local development: `http://localhost:8080`

> The default port is `8080`. Use this unless a different value for `PORT` was
> specified in the `.env` file.

## Endpoints

### POST /api/users

Creates a new user.

- **Request body**

  ```json
  {
    "email": "alice@example.com",
    "password": "your_password"
  }
  ```

- **Successful response**

  `201 Created`

  ```json
  {
    "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z",
    "email": "alice@example.com",
    "is_chirpy_red": false
  }
  ```

- **Error responses**
  - `400 Bad Request`
  - `500 Internal Server Error`

### PUT /api/users

Updates email and password of the current user. Requires a valid JWT access
token in the `Authorization` header.

- **Request header**

  ```http
  Authorization: Bearer <access_token>
  ```

- **Request body**

  ```json
  {
    "email": "alice2@example.com",
    "password": "new_password"
  }
  ```

- **Successful response**

  `200 OK`

  ```json
  {
    "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-02-01T12:00:00Z",
    "email": "alice2@example.com",
    "is_chirpy_red": false
  }
  ```

- **Error responses**
  - `401 Unauthorized`
  - `500 Internal Server Error`

### POST /api/login

Logs in a user.

- **Request body**

  ```json
  {
    "email": "your_email@example.com",
    "password": "your_password"
  }
  ```

- **Successful response**

  `200 OK`

  ```json
  {
      "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
      "created_at": "2026-01-01T12:00:00Z",
      "updated_at": "2026-01-01T12:00:00Z",
      "email": "alice@example.com",
      "is_chirpy_red": false,
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.access-token-part.signature-part",
      "refresh_token": "random-256-bit-string"
  }
  ```

  > For brevity, the example tokens above are truncated.

- **Error responses**
  - `401 Unauthorized`
  - `500 Internal Server Error`

### POST /api/refresh

Returns a new access token. Requires a valid refresh token in the
`Authorization` header.

- **Request header**

  ```http
  Authorization: Bearer <refresh_token>
  ```

- **Request body**
  
  No body.

- **Successful response**

  `200 OK`

  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.access-token-part.signature-part"
  }
  ```

  > For brevity, the example tokens above are truncated.

- **Error responses**
  - `401 Unauthorized`
  - `500 Internal Server Error`

### POST /api/revoke

Revokes a refresh token. Requires a valid refresh token in the `Authorization`
header. This refresh token will be the one to be revoked.

- **Request header**

  ```http
  Authorization: Bearer <refresh_token>
  ```

- **Request body**
  
  No body.

- **Successful response**

  `204 No Content` (no body)

- **Error responses**
  - `401 Unauthorized`
  - `500 Internal Server Error`

### GET /api/chirps

Returns a list of chirps.

- **Query parameters**
  - `author_id` - filters chirps by the author's user ID
  - `sort` - sorts the response list of chirps:
    - `asc` - ascending by create date (default)
    - `desc` - descending by create date

- **Request body**
  
  No body.

- **Successful response**
  
  `200 OK`

  ```json
  [
    {
      "id": "fdeb872f-c649-4e8e-a03e-41970d724f67",
      "created_at": "2026-01-15T15:48:01.817181Z",
      "updated_at": "2026-01-15T15:48:01.817181Z",
      "body": "Hello from Chirpy!",
      "user_id": "94b13964-b8ab-4027-a43a-7fb0f94cea6e"
    },
    {
      "id": "722c9d67-aa68-4287-97f7-fb2c8c7d7b0b",
      "created_at": "2026-01-15T15:48:01.830566Z",
      "updated_at": "2026-01-15T15:48:01.830566Z",
      "body": "Hello again!",
      "user_id": "94b13964-b8ab-4027-a43a-7fb0f94cea6e"
    },
    {
      "id": "28fe52eb-ca3c-4e96-9d8e-a419ae1d5c82",
      "created_at": "2026-01-15T15:48:01.843303Z",
      "updated_at": "2026-01-15T15:48:01.843303Z",
      "body": "Hello a third time!",
      "user_id": "94b13964-b8ab-4027-a43a-7fb0f94cea6e"
    }
  ]
  ```

- **Error responses**
  - `400 Bad Request`
  - `404 Not Found`
  - `500 Internal Server Error`

### GET /api/chirps/{id}

Returns a single chirp by ID.

- **Path parameter**
  - `id` - chirp ID

- **Request body**
  
  No body.

- **Successful response**
  
  `200 OK`

  ```json
  {
    "id": "4b3cab61-d91e-4024-9a42-1365e22953f5",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z",
    "body": "Hello from Chirpy!",
    "user_id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b"
  }
  ```

- **Error responses**
  - `400 Bad Request`
  - `404 Not Found`
  - `500 Internal Server Error`

### POST /api/chirps

Creates a new chirp. Requires a valid JWT access token in the `Authorization`
header.

- **Request header**

  ```http
  Authorization: Bearer <access_token>
  ```

- **Request body**

  ```json
  {
    "body": "Hello from Chirpy!"
  }
  ```

- **Successful response**
  
  `201 Created`

  ```json
  {
    "id": "4b3cab61-d91e-4024-9a42-1365e22953f5",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z",
    "body": "Hello from Chirpy!",
    "user_id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b"
  }
  ```

- **Error responses**
  - `400 Bad Request`
  - `401 Unauthorized`
  - `500 Internal Server Error`

### DELETE /api/chirps/{id}

Deletes a chirp by ID. Requires a valid JWT access token in the `Authorization`
header. Only the owner of a chirp may delete it.

- **Request header**

  ```http
  Authorization: Bearer <access_token>
  ```

- **Path parameters**
  - `id` - chirp ID

- **Request body**
  
  No body.

- **Successful response**
  
  `204 No Content` (no body)

- **Error responses**
  - `400 Bad Request`
  - `401 Unauthorized`
  - `403 Forbidden`
  - `404 Not Found`
  - `500 Internal Server Error`

## Error Format

All 4xx/5xx error responses use this format:

```json
{
  "error": "description of what went wrong."
}
```
