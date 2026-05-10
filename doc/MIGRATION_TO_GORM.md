# Architecture Migration and GORM Implementation

This document describes the migration from the original monolithic service to a SOLID-principles-based layered architecture using GORM for data persistence.

This document describes the migration from `database/sql` with `go-sql-driver/mysql` to GORM ORM library.

## Summary of Changes

### 1. Dependencies
**Added:**
- `gorm.io/gorm` - GORM ORM library
- `gorm.io/driver/mysql` - MySQL driver for GORM
- `gorm.io/driver/sqlite` - SQLite driver for GORM (testing only)

**Removed:**
- `github.com/DATA-DOG/go-sqlmock` - No longer needed for testing
- Direct usage of `database/sql` package

**Note:** `github.com/go-sql-driver/mysql` is still present as an indirect dependency (used by GORM's MySQL driver).

### 2. Code Changes

#### New Model File: `internal/guestbook/model.go`
Created a GORM model for the guestbook entry:
```go
type GuestbookEntry struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    Name      string    `gorm:"type:varchar(255);not null"`
    Message   string    `gorm:"type:text;not null"`
    CreatedAt int64     `gorm:"not null"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

#### Updated: `internal/guestbook/service.go`
**Before:**
```go
type Service struct {
    pb.UnimplementedGuestbookServiceServer
    db *sql.DB
}

func (s *Service) AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
    query := "INSERT INTO guestbook_entries (name, message, created_at) VALUES (?, ?, ?)"
    _, err := s.db.ExecContext(ctx, query, req.Name, req.Message, time.Now().Unix())
    // ...
}
```

**After:**
```go
type Service struct {
    pb.UnimplementedGuestbookServiceServer
    db *gorm.DB
}

func (s *Service) AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
    entry := &GuestbookEntry{
        Name:      req.Name,
        Message:   req.Message,
        CreatedAt: time.Now().Unix(),
    }
    if err := s.db.WithContext(ctx).Create(entry).Error; err != nil {
        // ...
    }
}
```

#### Updated: `cmd/server/main.go`
**Before:**
```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

db, err := sql.Open("mysql", dsn)
if err != nil {
    log.Fatalf("Failed to open database: %v", err)
}
defer db.Close()
```

**After:**
```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}

// Auto-migrate the schema
if err := db.AutoMigrate(&guestbook.GuestbookEntry{}); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
}

// Get underlying SQL DB for connection pooling
sqlDB, err := db.DB()
if err != nil {
    log.Fatalf("Failed to get database instance: %v", err)
}
defer sqlDB.Close()
```

#### Updated: `internal/guestbook/service_test.go`
**Before:**
- Used `go-sqlmock` to mock database operations
- Required setting up expectations for each query
- Tested against mocked responses

**After:**
- Uses in-memory SQLite database for testing
- Tests run against a real database
- More realistic integration testing
- No need to mock database calls

### 3. Benefits of GORM

1. **Cleaner Code**: No more raw SQL queries in the code
2. **Type Safety**: Compile-time checking of model fields
3. **Auto-Migration**: Database schema is automatically created/updated
4. **Better Testing**: Tests use real database instead of mocks
5. **Context Support**: Built-in context support for all operations
6. **Easier Maintenance**: ORM handles SQL generation and optimization
7. **Database Portability**: Easier to switch between different databases

### 4. Testing Improvements

**Before:**
- Used `sqlmock` to mock database calls
- Required setting up expectations for each test
- Didn't catch real database issues

**After:**
- Uses in-memory SQLite for testing
- Each test gets a fresh database
- Tests run against real database operations
- Catches more potential issues
- Still maintains 100% code coverage

### 5. Database Schema Management

GORM now handles schema management through auto-migration:
```go
db.AutoMigrate(&guestbook.GuestbookEntry{})
```

This automatically:
- Creates tables if they don't exist
- Adds missing columns
- Updates column types (with limitations)
- Creates indexes defined in the model

**Note:** For production, consider using proper migration tools like `golang-migrate` or GORM's migration features with version control.

### 6. Performance Considerations

- GORM adds a small overhead compared to raw SQL
- For this application, the overhead is negligible
- GORM provides connection pooling through the underlying `sql.DB`
- Can access the underlying `sql.DB` for advanced configuration:
  ```go
  sqlDB, err := db.DB()
  sqlDB.SetMaxOpenConns(100)
  sqlDB.SetMaxIdleConns(10)
  ```

### 7. Migration Checklist

- [x] Install GORM and drivers
- [x] Create GORM models
- [x] Update service to use GORM
- [x] Update main.go to use GORM
- [x] Rewrite tests to use in-memory SQLite
- [x] Update documentation
- [x] Remove old dependencies
- [x] Test all functionality
- [x] Verify build succeeds

## Running the Application

The application works exactly the same as before from a user perspective:

```bash
# Run the server
make run

# Run tests
go test -v ./internal/guestbook/...

# Check coverage
go test -cover ./internal/guestbook/...
```

## Rollback Plan

If needed to rollback:
1. Revert changes to `service.go`, `main.go`, and `service_test.go`
2. Remove `model.go`
3. Run `go get github.com/go-sql-driver/mysql@v1.9.3`
4. Run `go get github.com/DATA-DOG/go-sqlmock@v1.5.2`
5. Run `go mod tidy`

## Future Considerations

1. **Advanced Queries**: GORM supports complex queries, joins, and aggregations
2. **Hooks**: Can add before/after save hooks for validation or logging
3. **Soft Deletes**: Easy to implement soft deletes with GORM
4. **Pagination**: GORM has built-in pagination support
5. **Transactions**: GORM provides transaction support with rollback capabilities
6. **Database Migrations**: Consider using GORM Migrator or external tools for production migrations
