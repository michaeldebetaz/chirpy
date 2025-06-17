# Chirpy API

Chirpy is a microblogging REST API written in Go. This API enables users to manage accounts, post and retrieve "chirps" (short messages), and supports authentication, tokens, and admin metrics.

## Base URL

```
http://localhost:8080
```

---

## Authentication

Most endpoints require JWT authentication via the Authorization header:

```
Authorization: Bearer <token>
```

Some endpoints (like registration, login, and health check) do not require authentication.

---

## Endpoints

### Health Check

- **GET /api/healthz**
  - Returns API health status.
  - Response: `200 OK`, body: `"OK"`

---

### Chirps

- **GET /api/chirps**
  - List all chirps.
  - Optional query: `author_id` (filter by user).
  - Response: `200 OK`, JSON array of chirps.

- **GET /api/chirps/{chirpID}**
  - Get a specific chirp by ID.
  - Response: `200 OK`, JSON chirp object.
  - Errors: `404 Not Found` if not found.

- **POST /api/chirps**
  - Create a new chirp (requires JWT).
  - Request body:
    ```json
    {
      "body": "string"
    }
    ```
  - Response: `201 Created`, JSON chirp object.

- **DELETE /api/chirps/{chirpID}**
  - Delete a chirp (requires JWT). Only the author can delete.
  - Response: `204 No Content`.
  - Errors: `403 Forbidden` if not the author.

---

### Users

- **POST /api/users**
  - Register a new user.
  - Request body:
    ```json
    {
      "email": "user@example.com",
      "password": "password"
    }
    ```
  - Response: `201 Created`, user JSON.

- **PUT /api/users**
  - Update user email or password (requires JWT).
  - Request body:
    ```json
    {
      "email": "new@example.com",
      "password": "new_password"
    }
    ```
  - Response: `200 OK`, updated user JSON.

---

### Authentication

- **POST /api/login**
  - Login and receive access/refresh tokens.
  - Request body:
    ```json
    {
      "email": "user@example.com",
      "password": "password",
      "expires_in_seconds": 3600
    }
    ```
  - Response: `200 OK`, JSON with JWT and refresh token.

- **POST /api/refresh**
  - Refresh JWT using refresh token.
  - Requires Authorization: Bearer <refresh_token>
  - Response: `200 OK`, new JWT.

- **POST /api/revoke**
  - Revoke a refresh token.
  - Requires Authorization: Bearer <refresh_token>
  - Response: `204 No Content`.

---

### Admin & Metrics

- **GET /admin/metrics**
  - View metrics (hits).
  - Response: `200 OK`, HTML.

- **POST /admin/reset**
  - Reset all users and metrics (dev mode only).
  - Response: `200 OK`.

---

### Webhooks

- **POST /api/polka/webhooks**
  - Accepts Polka webhooks for user upgrades.
  - Requires `Polka-Key` header.

---

## Error Responses

- Errors are returned as JSON:
  ```json
  {"error": "message"}
  ```

- Common status codes: `400`, `401`, `403`, `404`, `500`

---

## Running Locally

1. Clone the repo.
2. Build and run with Go:
   ```
   go run main.go
   ```
3. The API will be available at `http://localhost:8080`
