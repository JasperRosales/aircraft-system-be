# User Service Documentation

## Overview

The User Service handles user authentication and CRUD operations using JWT tokens with HTTP-only cookies.

## Architecture

```
cmd/api/main.go
    ↓
internal/routers/user_router.go
    ↓
internal/controller/user_controller.go
    ↓
internal/service/user_service.go
    ↓
internal/repository/user_repo.go
    ↓
Database (Supabase/PostgreSQL)
```

## API Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|----------------|-------------|
| POST | `/api/users/register` | No | Register a new user |
| POST | `/api/users/login` | No | Login and receive auth cookie |
| POST | `/api/users/logout` | No | Clear auth cookie |
| GET | `/api/users/me` | Yes | Get current authenticated user |
| GET | `/api/users/:id` | Yes | Get user by ID |
| GET | `/api/users` | Yes | Get all users |
| PUT | `/api/users/:id` | Yes | Update user |
| DELETE | `/api/users/:id` | Yes | Delete user |

## Authentication Middleware

### AuthMiddleware

Validates JWT tokens from either:
1. HTTP-only `auth_token` cookie (preferred)
2. `Authorization: Bearer <token>` header

**Response on missing/invalid token:**
```json
{
  "error": "authentication required"
}
```

**Response on expired token:**
```json
{
  "error": "token has expired"
}
```

**Response on invalid token:**
```json
{
  "error": "invalid token"
}
```

### OptionalAuthMiddleware

Similar to AuthMiddleware but doesn't reject requests without tokens. Useful for endpoints that change behavior based on authentication but don't require it.

### RoleMiddleware

Checks if the authenticated user has the required role:
```go
protected.Use(middleware.RoleMiddleware("admin"))
```

**Response on insufficient permissions:**
```json
{
  "error": "insufficient permissions"
}
```

## Request/Response Models

### Register Request
```json
{
  "name": "string (required, 2-255 chars)",
  "password": "string (required, min 6 chars)",
  "role": "string (optional, 'user', 'mechanic', or 'admin', default: 'user')"
}
```

### Login Request
```json
{
  "name": "string (required)",
  "password": "string (required)"
}
```

### Get Current User Request
**Endpoint:** `GET /api/users/me`

**Authentication:** Required (JWT token via cookie or Bearer header)

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "testuser",
  "role": "user",
  "created_at": "2026-02-10T12:34:56Z"
}
```

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `404 Not Found`: User not found

### Update Request
```json
{
  "name": "string (optional, 2-255 chars)",
  "password": "string (optional, min 6 chars)",
  "role": "string (optional, 'user' or 'admin')"
}
```

### User Response
```json
{
  "id": 1,
  "name": "testuser",
  "role": "user",
  "created_at": "2026-02-10T12:34:56Z",
  "updated_at": "2026-02-10T12:34:56Z"
}
```

## JWT Token (HTTP-only Cookie)

- **Cookie Name:** `auth_token`
- **HttpOnly:** Yes (not accessible via JavaScript)
- **Secure:** Yes (HTTPS only)
- **SameSite:** Strict
- **Default Expiry:** 24 hours (configurable via `TOKEN_EXP` env var)

### Token Claims
```go
type JWTClaims struct {
    UserID int64  `json:"user_id"`
    Name   string `json:"name"`
    Role   string `json:"role"`
}
```

## Password Security

- Passwords are hashed using **bcrypt** with default cost (10)
- Passwords are never returned in API responses
- Password is excluded from JSON serialization using `json:"-"` tag

## User Roles

| Role | Description |
|------|-------------|
| `user` | Default role for regular users |
| `mechanic` | Maintenance personnel with parts management access |
| `admin` | Administrator with elevated privileges |

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | - | Supabase connection string |
| `SECRET` | No | Default key | JWT signing secret |
| `TOKEN_EXP` | No | `24` | Token expiry in hours |
| `PORT` | No | `8080` | Server port |

## Testing with curl

```bash
# 1. Register (public)
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"name":"testuser","password":"password123"}'

# 2. Login (public) - gets cookie
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"name":"testuser","password":"password123"}' \
  -c cookies.txt

# 3. Get all users (protected) - requires cookie
curl -X GET http://localhost:8080/api/users -b cookies.txt

# 4. Get current user (protected)
curl -X GET http://localhost:8080/api/users/me -b cookies.txt

# 5. Get user by ID (protected)
curl -X GET http://localhost:8080/api/users/1 -b cookies.txt

# 6. Update user (protected)
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"newname"}' \
  -b cookies.txt

# 7. Delete user (protected)
curl -X DELETE http://localhost:8080/api/users/1 -b cookies.txt

# 8. Logout (clears cookie)
curl -X POST http://localhost:8080/api/users/logout -c cookies.txt

# Testing without auth (will fail for protected routes)
curl -X GET http://localhost:8080/api/users
# Response: {"error":"authentication required"}
```

## Using Bearer Token Header

```bash
# Login and save token
TOKEN=$(curl -s -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"name":"testuser","password":"password123"}' | jq -r '.token')

# Use Bearer token instead of cookie
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"
```


## Response Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (user exists) |
| 500 | Internal Server Error |

