# Authentication App

A simple authentication application with JWT tokens and file upload functionality built with Go and Fiber.

## Features

- User registration and login
- JWT-based authentication
- File upload with authentication
- Token revocation
- Database migrations
- Swagger API documentation
- Health check endpoints

## Prerequisites

- Go 1.21+
- PostgreSQL
- Docker & Docker Compose (optional)

## Installation

```bash
go mod download
```

## Running with Docker Compose

```bash
# Start the application with database
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop the application
docker-compose down
```

## Running Locally

### 1. Start Database

```bash
# Using Docker
docker run --name postgres-auth -e POSTGRES_DB=authentication-app -e POSTGRES_USER=user123 -e POSTGRES_PASSWORD=pass123 -p 5432:5432 -d postgres:15
```

### 2. Run Migrations

```bash
# Install migrate package
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrate
migrate -path migrations -database "postgresql://user123:pass123@localhost:5432/authentication-app?sslmode=disable" up
```

### 3. Start Application

```bash
# Production mode
go run cmd/app/main.go
```

## API Documentation

- **Swagger UI**: `http://localhost:8080/api/swagger`
- **Web Interface**: `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user
- `POST /auth/revoke` - Revoke JWT token (requires authentication)

### File Upload

- `POST /files/upload` - Upload file (requires authentication)

### Health Check

- `GET /api/readiness` - Readiness probe
- `GET /api/liveness` - Liveness probe

## Testing File Upload

### 1. Register/Login via Web Interface

Open `http://localhost:8080` and register or login to get a JWT token.

### 2. Upload via Web Interface

Use the file upload form after logging in.

### 3. Upload via API (curl)

```bash
# First, login to get a token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your_username","password":"your_password"}'

# Use the token from response to upload a file
curl -X POST http://localhost:8080/files/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -F "file=@/path/to/your/file.jpg"
```

## Development

The application uses Air for hot reloading during development. Configuration is in `.air.toml`.

```bash
# Install Air
go install github.com/air-verse/air@latest

# Start development server
air
```

## Project Structure

```
├── cmd/                                        # Application entry points
│   └── app/
│       └── main.go                             # Main application entry point - initializes server and dependencies
├── config/                                     # Configuration management
│   └── config.go                               # Application configuration settings (database, server, JWT settings)
├── docs/                                       # API documentation files
│   ├── docs.go                                 # Generated Swagger documentation code
│   ├── swagger.json                            # Swagger API specification in JSON format
│   └── swagger.yaml                            # Swagger API specification in YAML format
├── internal/                                   # Private application code (not importable by other projects)
│   ├── controllers/                            # HTTP request handlers (Controller layer)
│   │   ├── auth.controller.go                  # Authentication endpoints (register, login, revoke token)
│   │   ├── file.controller.go                  # File upload and management endpoints
│   │   └── monitor.controller.go               # Health check and monitoring endpoints
│   ├── DTOS/                                   # Data Transfer Objects for API requests/responses
│   │   └── auth_dto.go                         # Authentication-related DTOs (login, register requests)
│   ├── middleware/                             # HTTP middleware functions
│   │   └── jwt.go                              # JWT authentication middleware for protecting routes
│   ├── models/                                 # Database models and business entities
│   │   ├── file_upload.go                      # File upload metadata model
│   │   ├── revoked_token.go                    # Revoked JWT tokens model for security
│   │   └── user.go                             # User model with authentication fields
│   ├── server/                                 # Server setup and routing configuration
│   │   ├── handlers.go                         # Route handlers registration and middleware setup
│   │   └── server.go                           # Fiber server initialization and configuration
├── migrations/                                 # Database schema migrations
│   ├── 001_create_users_table.up.sql           # Creates users table with authentication fields
│   ├── 002_create_revoked_tokens_table.up.sql  # Creates table for tracking revoked JWT tokens
│   ├── 003_create_file_uploads_table.up.sql    # Creates table for file upload metadata
├── pkg/                                        # Public packages (reusable across projects)
│   ├── database/                               # Database connection and management
│   │   └── postgresql.go                       # PostgreSQL connection setup and configuration
│   └── utils/                                  # Utility functions and helpers
│       ├── auth.go                             # Authentication utilities (password hashing, validation)
│       ├── file.go                             # File handling utilities (validation, storage)
│       └── token.go                            # JWT token generation, validation, and management
├── .air.toml                                   # Hot reload configuration for development
├── .env                                        # Environment variables (database credentials, JWT secrets)
├── .gitignore                                  # Git ignore rules for excluding sensitive files
├── docker-compose.yml                          # Docker compose setup for app and PostgreSQL
├── Dockerfile                                  # Docker image configuration for the application
├── go.mod                                      # Go module dependencies declaration
├── go.sum                                      # Go module dependency checksums
└── README.md                                   # Project documentation and setup instructions
```
