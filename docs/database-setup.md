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

### planes
```sql
CREATE TABLE planes (
    id SERIAL PRIMARY KEY,
    tail_number VARCHAR(50) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### plane_parts
```sql
CREATE TABLE plane_parts (
    id SERIAL PRIMARY KEY,
    plane_id INTEGER NOT NULL REFERENCES planes(id) ON DELETE CASCADE,
    part_name VARCHAR(255) NOT NULL,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(150) NOT NULL,
    usage_hours NUMERIC(10,2) DEFAULT 0,
    usage_limit_hours NUMERIC(10,2) NOT NULL,
    usage_percent NUMERIC GENERATED ALWAYS AS 
        ((usage_hours / usage_limit_hours) * 100) STORED,
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
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

