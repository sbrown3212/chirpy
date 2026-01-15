# Chirpy API

A simple REST API for creating and viewing short text posts ("chirps").

## Learning Project Disclaimer

This project was built as part of a guided learning project from
[Boot.dev](https://boot.dev) and is not intended for production use.

## Features

- Create and update users
- Create new chirps
- List all chirps or get chirp by id
- Uses access and refresh tokens for authentication

## Setup

### Requirements

- Go 1.23+
- PostgreSQL 15+
- `psql` command line client
- Goose

### Create Database

Create a database for Chirpy:

```
createdb chirpy
```

### Clone Repository

```
git clone https://github.com/sbrown3212/chirpy.git
```

### Environment Variables

There should be a file named `.env.example` at the root of the project. This is an
example of what the Chirpy server expects the `.env` file to look like, as well
as explanations of what values should be used.

Rename this `.env.example` file to `.env`, or copy it to a file named `.env`.
Then, update the variables as needed.

> Be sure not to commit your .env file to version control.

### Run Database Migrations

Load environment variables for CLI tools (like goose):

```
source .env
```

Run migrations (uses $DB_URL from .env):

```
goose postgres "$DB_URL" up
```

### Run Locally

Navigate to project root and run:

```
go run .
```

## Base URL

- `http://localhost:8080` (unless PORT environment variable was changed)

## Endpoints

### POST /api/users

Creates a new user.

- Request Body

  ```
  {
    "email": "alice@example.com",
    "password": "your_password"
  }
  ```

- Successful Response (`201 Created`):

  ```
  {
    "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z",
    "email": "alice@example.com",
    "is_chirpy_red": false
  }
  ```

- Error Responses
  - `400 Bad Request`
  - `500 Internal Server Error`

### PUT /api/users

Updates email and password of the current user. Requires a valid JWT access
token.

- Request Header

  ```
  Authorization: Bearer <access_token>
  ```

- Request Body

  ```
  {
    "email": "alice2@example.com",
    "password": "new_password"
  }
  ```

- Successful Response (`200 OK`):

  ```
  {
    "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z",
    "email": "alice@example.com",
    "is_chirpy_red": false
  }
  ```

- Error Responses
  - `401 Unauthorized`
  - `500 Internal Server Error`

### POST /api/login

Logs in a user. Response contains user information and JWT access and refresh
tokens.

- Request Header

  ```
  Authorization: Bearer <access_token>
  ```

- Request body:

  ```
  {
    "email": "your_email@example.com",
    "password": "your_password"
  }
  ```

- Successful Response (`200 OK`):

  ```
  {
      "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
      "created_at": "2026-01-01T12:00:00Z",
      "updated_at": "2026-01-01T12:00:00Z",
      "email": "<alice@example.com>",
      "is_chirpy_red": false,
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.access-token-part.signature-part",
      "refresh_token": "random-256-bit-string"
  }
  ```

  > For brevity, the example tokens above are truncated.

- Error Responses:
  - `401 Unauthorized`
  - `500 Internal Server Error`

### POST /api/refresh

Returns a new access token. Requires a valid refresh token in the Authorization
header.

- Request Header

  ```
  Authorization: Bearer <refresh_token>
  ```

- Request Body: none

- Successful Response (`200 OK`):

  ```
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.access-token-part.signature-part"
  }
  ```

  > For brevity, the example tokens above are truncated.

- Error Responses:
  - `401 Unauthorized`
  - `500 Internal Server Error`

### POST /api/revoke

Revokes a refresh token. Requires a valid refresh token in the Authorization
header.

- Request Header

  ```
  Authorization: Bearer <refresh_token>
  ```

- Request Body: none

- Successful Response (`204 No Content`): no body

- Error Responses:
  - `401 Unauthorized`
  - `500 Internal Server Error`

### GET /api/chirps

Returns a list of chirps.

- Query Parameters
  - `author_id` - filter chirps by author's user id
  - `sort` - sorts response list of chirps, and accepts the following values:
    - `asc` - ascending by create date (default)
    - `desc` - descending by create date

- Request Body: none

- Successful Response (`200 OK`):

  ```
  [
    {
      "id": "fdeb872f-c649-4e8e-a03e-41970d724f67",
      "created_at": "2026-01-15T15:48:01.817181Z",
      "udpated_at": "2026-01-15T15:48:01.817181Z",
      "body": "Hello from Chirpy!",
      "user_id": "94b13964-b8ab-4027-a43a-7fb0f94cea6e"
    },
    {
      "id": "722c9d67-aa68-4287-97f7-fb2c8c7d7b0b",
      "created_at": "2026-01-15T15:48:01.830566Z",
      "udpated_at": "2026-01-15T15:48:01.830566Z",
      "body": "Hello again!",
      "user_id": "94b13964-b8ab-4027-a43a-7fb0f94cea6e"
    },
    {
      "id": "28fe52eb-ca3c-4e96-9d8e-a419ae1d5c82",
      "created_at": "2026-01-15T15:48:01.843303Z",
      "udpated_at": "2026-01-15T15:48:01.843303Z",
      "body": "Hello a third time!",
      "user_id": "94b13964-b8ab-4027-a43a-7fb0f94cea6e"
    }
  ]
  ```

- Error Responses:
  - `400 Bad Request`
  - `404 Not Found`
  - `500 Internal Server Error`

### GET /api/chirps/{id}

Returns a single chirp by ID.

- Path Parameter:
  - `id` - chirp ID

- Request Body: none

- Successful Response (`200 OK`):

  ```
  {
    "id": "4b3cab61-d91e-4024-9a42-1365e22953f5",
    "created_at": "2026-01-01T12:00:00Z",
    "udpated_at": "2026-01-01T12:00:00Z",
    "body": "Hello from Chirpy!",
    "user_id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b"
  }
  ```

- Error Responses:
  - `400 Bad Request`
  - `404 Not Found`
  - `500 Internal Server Error`

### POST /api/chirps

Creates a new chirp. Requires a valid JWT access token in Authorization header.

- Request body:

  ```
  {
    "body": "Hello from Chirpy!"
  }
  ```

- Successful Response (`201 Created`)

  ```
  {
    "id": "4b3cab61-d91e-4024-9a42-1365e22953f5",
    "created_at": "2026-01-01T12:00:00Z",
    "udpated_at": "2026-01-01T12:00:00Z",
    "body": "Hello from Chirpy!",
    "user_id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b"
  }
  ```

- Error Responses:
  - `400 Bad Request`
  - `401 Unauthorized`
  - `500 Internal Server Error`

### DELETE /chirps/{id}

Deletes a chirp by ID. Requires a valid JWT access token in Authorization
header.

- Path Parameters:
  - `id` - chirp ID

- Request Body: none

- Successful Response (`204 No Content`): no body

- Error Responses:
  - `400 Bad Request`
  - `401 Unauthorized`
  - `403 Forbidden`
  - `404 Not Found`
  - `500 Internal Server Error`

## Error Format

All error responses use:

```
{
  "error": "description of what whent wrong."
}
```
