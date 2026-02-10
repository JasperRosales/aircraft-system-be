# Supabase Setup with Goose and GORM

This project uses Supabase (PostgreSQL) with Goose for migrations and GORM as the ORM.

## Database Sche
```

## Environment Variables

Set these in your `.env` file:
```bash
# Supabase Connection String
DATABASE_URL="postgresql://user:password@host:5432/dbname?sslmode=require&search_path=public"

# JWT Configuration
SECRET="your-jwt-secret-key"
TOKEN_EXP="3"  # Token expiry in hours
```

## Goose Migrations

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

## GORM Setup

### Database Initialization

The project initializes GORM in `cmd/api/main.go`:

```go
func initDatabase() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        return nil, nil
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    return db, nil
}
```

### GORM Model Definition

```go
type User struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Password  string    `json:"-" db:"password"`
    Role      string    `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**Note:** The `json:"-"` tag prevents password from being serialized in JSON responses.

## Connecting to Supabase

1. Get your Supabase connection string from:
   - Dashboard → Settings → Database → Connection string
   - Use the "Transaction" pooler URL for development

2. Enable SSL requirement:
   ```bash
   ?sslmode=require
   ```

3. Set the `DATABASE_URL` environment variable.

## Troubleshooting

### Connection Issues
- Ensure IP whitelist includes your current IP
- Check that the password doesn't contain special characters that need URL encoding

### Migration Errors
- Run `goose status` to check current migration state
- Use `goose redo` to rollback and re-apply the last migration

