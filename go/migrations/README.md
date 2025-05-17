# Database Migrations

This directory contains database migrations for the application using [Goose](https://github.com/pressly/goose).

## Prerequisites

- PostgreSQL 12+
- `uuid-ossp` extension
- `pgvector` extension (for vector embeddings)

## Setup

1. Install the required PostgreSQL extensions:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgvector";
```

## Running Migrations

### Using Goose CLI

1. Install Goose:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

2. Run migrations:
```bash
# Up
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable" \
  goose -dir migrations/postgres postgres "$DATABASE_URL" up

# Down (rollback one migration)
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable" \
  goose -dir migrations/postgres postgres "$DATABASE_URL" down

# Status
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable" \
  goose -dir migrations/postgres postgres "$DATABASE_URL" status
```

## Creating New Migrations

To create a new migration:

```bash
go run github.com/pressly/goose/v3/cmd/goose -dir migrations/postgres postgres "$DATABASE_URL" create add_something sql
```

## Database Schema

### Stores
- `id`: Auto-incrementing integer (SERIAL)
- `type`: Type of the store (e.g., 'github')
- `repo`: Repository identifier (for GitHub stores)
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update

### Documents
- `id`: UUID (auto-generated)
- `store_id`: Foreign key to stores.id
- `path`: Path to the document
- `content`: Document content
- `embedding`: Vector embedding of the document (pgvector)
- `tags`: JSONB array of tags
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update

### Memories
- `id`: UUID (auto-generated)
- `path`: Path to the memory
- `content`: Memory content
- `embedding`: Vector embedding of the memory (pgvector)
- `tags`: JSONB array of tags
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update
