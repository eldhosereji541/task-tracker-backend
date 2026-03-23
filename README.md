# Task Tracker Backend

A GraphQL API for managing tasks, built with Go. Supports user auth (register/login with JWT) and full CRUD on tasks — filter by status, priority, or search term.

## Stack

- **Go** — language
- **Echo** — HTTP server
- **gqlgen** — GraphQL code generation
- **JWT** — auth tokens
- **bcrypt** — password hashing
- In-memory store (no database, data resets on restart)

## Getting Started

**1. Clone and install dependencies**

```bash
git clone https://github.com/eldhosereji541/task-tracker-backend
cd task-tracker-backend
go mod download
```

**2. Set up environment variables**

```bash
cp .env.example .env
```

Edit `.env` and fill in your values:

```env
JWT_SECRET=your-secret-key-at-least-32-characters-long
PORT=8080
```

> JWT_SECRET must be at least 32 characters or the server won't start.

**3. Run the server**

```bash
go run ./cmd/server
```

Open [http://localhost:8080](http://localhost:8080) for the GraphQL playground.


### Task Status & Priority

| Status | Values |
|--------|--------|
| TaskStatus | `TODO`, `IN_PROGRESS`, `COMPLETED` |
| Priority | `LOW`, `MEDIUM`, `HIGH` |

## Project Structure

```
.
├── cmd/server/         # entry point
├── internal/
│   ├── auth/           # JWT, password hashing, middleware
│   ├── graph/          # gqlgen resolvers and schema
│   ├── model/          # domain models
│   ├── store/          # in-memory store
│   └── repository/     # repository interfaces
```

## Running Tests

```bash
go test ./...
```
