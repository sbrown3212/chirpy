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

A file named `.env.example` should be in the root of the project. This is an
example of what the Chirpy server expects the `.env` file to look like, as well
as explanations of what values should be used.

Rename this `.env.example` file to `.env`, or copy it to a file named `.env`.
Then, update the variables as needed.

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

### POST /users

Creates a new user.

#### Request Body

```
{
  "email": "your_email@example.com",
  "password": "your_password"
}
```

#### Successful Response (`201 Created`)

```
{
  "id": "6f23ecb5-41d5-426a-bacc-1d13e96dab2b",
  "created_at": "2026-01-01T12:00:00Z",
  "updated_at": "2026-01-01T12:00:00Z",
  "email": "alice@example.com",
  "is_chirpy_red": false
}
```

#### Error Responses

- `400 Bad Request`
- `500 Internal Server Error`

### PUT /users

Updates email and password of the current user. Requires a valid JWT access
token.

#### Request Header

```
Authorization: Bearer <access_token>
```

#### Request body

```
{
  "email": "your_existing_or_new_email@example.com",
  "password": "your_existing_or_new_password"
}
```

- Response (200):

```
TODO
```

- TODO: failed responses

### POST /login

Logs in a user.

- Request body:

```
{
  "email": "your_email@example.com",
  "password": "your_password"
}
```

- Response (200)

```
TODO
```

- TODO: failed responses

### POST /refresh

Returns a new access token.

- TODO: request body (not required???)

- TODO: response bodies

### POST /revoke

Revokes a refresh token.

- TODO: request body (not required???)

- TODO: response bodies

### GET /chirps

Returns a list of chirps.

- Query Parameters
  - `author_id` - filter chirps by author id
  - `sort` - sorts response list of chirps, and accepts the following values:
    - `asc` - ascending by create date
    - `desc` - descending by create date

- TODO: response bodies

### GET /chirps/{id}

Returns a single chirp by ID.

- Path Parameter:
  - `id` - chirp ID

- TODO: response bodies

### POST /chirps

Creates a new chirp.

- Request body:

```
{
TODO
}
```

- TODO: response bodies

### DELETE /chirps/{id}

Deletes a chirp by ID.

- Path Parameters:
  - `id` - chirp ID

- TODO: response bodies

## Error Format

All error responses use:

```
{
  "error": "description of what whent wrong."
}
```
