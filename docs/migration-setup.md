# Migration Setup 

This project uses goose for migration and manipulation of database using the cli version. 

### Creating a New Migration

```bash
# Generate a new migration file
goose create <migration_name> sql
```

### Running Migrations

```bash
# Apply all migrations
goose -dir database/migrations postgres "$DATABASE_URL" up

# Rollback one migration
goose -dir database/migrations postgres "$DATABASE_URL" down

# Rollback all migrations
goose -dir database/migrations postgres "$DATABASE_URL" reset

# Check migration status
goose -dir database/migrations postgres "$DATABASE_URL" status

# Create table (without migrations)
goose -dir database/migrations postgres "$DATABASE_URL" create <name> sql
```

### Migration File Format

Example migration file (`database/migrations/YYYYMMDDHHMMSS_<name>.sql`):
```sql
-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
```