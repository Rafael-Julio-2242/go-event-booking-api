# Go Event Booking REST API

A simple event booking REST API built with Go, Gin, SQLite, and JWT authentication. It supports user signup/login, CRUD for events, and event registrations.

- **Framework**: Gin (`github.com/gin-gonic/gin`)
- **Database**: SQLite (`gorm.io/driver/sqlite`)
- **ORM:** Gorm (`github.com/go-gorm/gorm`)
- **Auth**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Password Hashing**: bcrypt (`golang.org/x/crypto/bcrypt`)
- **Go Version**: declared in `go.mod` (`go 1.24.x`)

---

## Project Structure

```
.
├── api.db                         # SQLite database file (auto-created)
├── cmd/
│   └── api/
│       └── main.go                # App entrypoint (server startup)
├── internal/
│   ├── auth/
│   │   ├── handler.go             # Auth middleware (Bearer token)
│   │   └── service.go             # JWT generation/verification and password hashing
│   ├── event/
│   │   └── handler.go             # Event handlers
│   ├── registration/
│   │   └── handler.go             # Registration handlers (register/unregister)
│   ├── routes/
│   │   └── router.go              # Routes wiring
│   └── user/
│       └── handler.go             # Auth handlers (signup/login)
└── pkg/
    ├── db/
    │   └── gorm.go                # DB initialization
    └── models/
        ├── event.go               # Event model and DB operations
        ├── registration.go        # Registration model and DB operations
        └── user.go                # User model and DB operations
```

---

## How it works

- The server starts from `cmd/api/main.go` on port `:8080`, calls `db.InitDB()` which:
  - Opens/creates `api.db` (SQLite)
  - Runs auto-migrations for `users`, `events`, and `registrations`
- Routes are registered in `internal/routes/router.go`.
  - Public: `GET /events`, `GET /events/:id`, `POST /signup`, `POST /login`
  - Authenticated: `POST /events`, `PUT /events/:id`, `DELETE /events/:id`, `POST /events/:id/register`, `DELETE /events/:id/register`
- Auth middleware (`internal/auth/handler.go`) expects an `Authorization` header with the `Bearer` scheme.
- Passwords are hashed via bcrypt in `internal/auth/service.go`.
- JWT tokens are generated/verified in `internal/auth/service.go` (HS256). Default secret is hardcoded as `supersecret`.

---

## Quickstart

### Prerequisites
- Go 1.24+

### Install dependencies
```
go mod tidy
```

### Run the server
```
go run ./...
```
The API will be available at `http://localhost:8080`.

A local SQLite DB file `api.db` will be created in the project root with the required tables.

---

## Configuration

- JWT Secret: defined in `internal/auth/service.go` as `const secretKey = "supersecret"`.
  - For production, change this value and consider loading from environment variables.
- Port: hardcoded in `cmd/api/main.go` as `:8080`.
- Database: SQLite file `api.db` in the project root.

---

## Data Models

### Users (`users`)
- `id INTEGER PRIMARY KEY AUTOINCREMENT`
- `email TEXT NOT NULL UNIQUE`
- `password TEXT NOT NULL` (bcrypt hashed)

### Events (`events`)
- `id INTEGER PRIMARY KEY AUTOINCREMENT`
- `name TEXT NOT NULL`
- `description TEXT NOT NULL`
- `location TEXT NOT NULL`
- `dateTime DATETIME NOT NULL`
- `userId INTEGER` (creator/owner)

### Registrations (`registrations`)
- `id INTEGER PRIMARY KEY AUTOINCREMENT`
- `eventId INTEGER`
- `userId INTEGER`

---

## Authentication

- Sign up and login endpoints are public.
- After logging in, you receive a JWT (`HS256`) with claims: `email`, `userId`, `exp`.
- To access protected endpoints, send the token in the `Authorization` header using the Bearer scheme.
  - Example: `Authorization: Bearer <token>`

---

## Request/Response Conventions

- Content-Type: `application/json`
- Date/time format for `Event.DateTime`: RFC3339 (example: `2025-09-28T12:00:00Z`)
- Errors are returned as JSON: `{"message": "..."}` with appropriate HTTP status codes.

---

## Endpoints

### Auth

#### POST /signup
Create a new user.

Request
```json
{
  "email": "user@example.com",
  "password": "strongpassword"
}
```

Success Response (201)
```json
{ "message": "User created successfully" }
```

Error Responses
- 200 with message on bad request body (as implemented)
- 500 on DB failure

---

#### POST /login
Authenticate an existing user and return a JWT.

Request
```json
{
  "email": "user@example.com",
  "password": "strongpassword"
}
```

Success Response (200)
```json
{
  "message": "Logged in!",
  "token": "<jwt-token>"
}
```

Error Responses
- 200 with message on bad request body (as implemented)
- 401 on invalid credentials
- 500 on token generation failure

Curl example
```
curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"strongpassword"}'
```

---

### Events (Public)

#### GET /events
List all events.

Success Response (200)
```json
[
  {
    "ID": 1,
    "Name": "Go Conference",
    "Description": "All about Go",
    "Location": "SF",
    "dateTime": "2025-09-28T12:00:00Z",
    "UserID": 2
  }
]
```

---

#### GET /events/:id
Fetch a single event by ID.

Success Response (200)
```json
{
  "ID": 1,
  "Name": "Go Conference",
  "Description": "All about Go",
  "Location": "SF",
  "dateTime": "2025-09-28T12:00:00Z",
  "UserID": 2
}
```

Error Responses
- 400 on invalid ID
- 500 if not found or query error

---

### Events (Authenticated)
All endpoints below require the `Authorization` header with the raw JWT token string.

#### POST /events
Create a new event. The creator becomes the owner (`userId`).

Headers
```
Authorization: <token>
Content-Type: application/json
```

Request
```json
{
  "Name": "Go Conference",
  "Description": "All about Go",
  "Location": "SF",
  "dateTime": "2025-09-28T12:00:00Z"
}
```

Success Response (201)
```json
{
  "message": "Event created!",
  "event": {
    "ID": 1,
    "Name": "Go Conference",
    "Description": "All about Go",
    "Location": "SF",
    "dateTime": "2025-09-28T12:00:00Z",
    "UserID": 2
  }
}
```

Error Responses
- 400 on invalid body
- 500 on DB errors

---

#### PUT /events/:id
Update an event. Only the owner can update.

Headers
```
Authorization: <token>
Content-Type: application/json
```

Request (any updatable fields)
```json
{
  "Name": "Updated Name",
  "Description": "Updated desc",
  "Location": "NYC",
  "dateTime": "2025-10-01T16:00:00Z"
}
```

Success Response (200)
```json
{
  "message": "event updated successfully",
  "updated_event": {
    "ID": 1,
    "Name": "Updated Name",
    "Description": "Updated desc",
    "Location": "NYC",
    "dateTime": "2025-10-01T16:00:00Z",
    "UserID": 2
  }
}
```

Error Responses
- 400 on invalid ID/body
- 401 if not owner
- 500 on DB errors

---

#### DELETE /events/:id
Delete an event. Only the owner can delete.

Headers
```
Authorization: <token>
```

Success Response (200)
```json
{ "message": "Event Deleted successfully" }
```

Error Responses
- 400 on invalid ID
- 401 if not owner
- 500 on DB errors

---

### Event Registration (Authenticated)

#### POST /events/:id/register
Register the current user for an event.

Headers
```
Authorization: <token>
```

Success Response (201)
```json
{ "message": "User Registered successfully" }
```

Error Responses
- 400 if already registered
- 400 on missing/invalid event ID
- 500 on DB errors

---

#### DELETE /events/:id/register
Cancel the current user registration for an event.

Headers
```
Authorization: <token>
```

Success Response (201)
```json
{ "message": "User Unregistered successfully" }
```

Error Responses
- 400 if user was not registered
- 400 on missing/invalid event ID
- 500 on DB errors

---

## Curl Examples

Signup
```
curl -s -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"strongpassword"}'
```

Login (capture token)
```
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"strongpassword"}' | jq -r .token)
```

Create Event
```
curl -s -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"Name":"Go Conference","Description":"All about Go","Location":"SF","dateTime":"2025-09-28T12:00:00Z"}'
```

List Events
```
curl -s http://localhost:8080/events | jq
```

Register for Event (id=1)
```
curl -s -X POST http://localhost:8080/events/1/register \
  -H "Authorization: Bearer $TOKEN"
```

Unregister from Event (id=1)
```
curl -s -X DELETE http://localhost:8080/events/1/register \
  -H "Authorization: Bearer $TOKEN"
```
