## Install

```

go mod download

go install github.com/air-verse/air@latest

```

## Run

```

go run cmd/app/main.go

```

### Dev mode

```

air

```

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management.

### Installation

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Apply Migrations

```bash
# Apply all pending migrations
migrate -path migrations -database "postgresql://user123:pass123@localhost:5432/authentication-app?sslmode=disable" up

# Apply specific number of migrations
migrate -path migrations -database "postgresql://user123:pass123@localhost:5432/authentication-app?sslmode=disable" up 2
```
