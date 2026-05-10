# Testing the Guestbook Application

This document describes the testing strategy and how to run tests for the refactored guestbook application following SOLID principles.

## Test Structure

The application uses a comprehensive layered testing approach:

- **Unit Tests**: Test individual components in isolation using mocks
- **Integration Tests**: Test complete flows with real dependencies  
- **End-to-End Tests**: Test the complete application via gRPC

## Architecture Testing

### Controller Layer Tests (`internal/controller/`)
- **Unit Tests**: Test request/response handling with mock services
- **Integration Tests**: Test full flows with real database
- Tests cover all gRPC endpoints and error scenarios

### Service Layer Tests (`internal/service/`)
- **Unit Tests**: Test business logic with mock repositories and validators
- **Validation Tests**: Test input validation rules
- Tests cover success paths and error conditions

### Repository Layer Tests (`internal/repository/`)
- **Unit Tests**: Test data access with in-memory SQLite database
- Tests cover CRUD operations and query scenarios

## Running Tests

### All Tests
```bash
go test ./...
```

### Internal Package Tests
```bash
go test ./internal/...
```

### Specific Package Tests

#### Controller Tests
```bash
go test ./internal/controller/...
```

#### Service Tests  
```bash
go test ./internal/service/...
```

#### Repository Tests
```bash
go test ./internal/repository/...
```

### Integration Tests
```bash
go test ./test/...
```

### With Verbose Output
```bash
go test -v ./internal/...
```

### With Coverage
```bash
go test -cover ./internal/...
```

### Coverage Profile
```bash
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

### Specific Test
```bash
go test -v -run TestGuestbookController_AddEntry ./internal/controller/...
```

### Race Detection
```bash
go test -v -race -coverprofile=coverage.out ./internal/...
```

## Test Coverage

### Controller Layer
- Request/response mapping
- Error handling and validation responses
- Service integration

### Service Layer  
- Business logic validation
- Error propagation
- Repository integration

### Repository Layer
- Data persistence and retrieval
- Query operations and pagination
- Database transaction handling

## Test Patterns

### Dependency Injection Testing
Tests use dependency injection for easy mocking:

```go
// Controller tests use mock services
mockService := &MockGuestbookService{}
controller := NewGuestbookController(mockService)

// Service tests use mock repositories
mockRepo := &MockGuestbookRepository{}
service := NewGuestbookService(mockRepo, validator)
```

### Database Testing
Repository tests use in-memory SQLite for realistic testing:

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to create test database: %v", err)
    }
    return db
}
```

### Integration Testing
Full-stack tests verify complete functionality:

```go
// Creates real database and full dependency chain
controllerFactory := controller.NewGuestbookControllerFactory(db)
ctrl := controllerFactory.CreateController()
```

## Key Test Cases

### Controller Tests
- `TestNewGuestbookController`: Constructor validation
- `TestGuestbookController_AddEntry`: Request handling
- `TestGuestbookController_ListEntries`: Response formatting
- `TestGuestbookController_Integration_*`: End-to-end flows

### Service Tests
- `TestGuestbookService_AddEntry_Success`: Business logic
- `TestGuestbookService_AddEntry_ValidationError`: Error handling
- `TestGuestbookService_ListEntries_*`: Query operations

### Repository Tests
- `TestGormGuestbookRepository_Create`: Data insertion
- `TestGormGuestbookRepository_FindWithPagination`: Data retrieval
- Context cancellation and error scenarios

### Validator Tests
- `TestDefaultValidator_ValidateAddEntryRequest`: Input validation
- `TestDefaultValidator_ValidateListEntriesRequest`: Pagination validation
- Edge cases and error conditions

## Benefits of This Testing Approach

### 1. **Isolation**
- Each layer tested independently
- Mock dependencies prevent cascading failures
- Fast unit test execution

### 2. **Confidence**
- Integration tests verify component interactions
- Real database tests catch SQL issues
- End-to-end tests validate user flows

### 3. **Maintainability**
- Clear test structure matches code architecture
- Easy to locate and update tests
- Consistent testing patterns

### 4. **Extensibility**
- Easy to add tests for new features
- Mock interfaces support new implementations
- Test utilities are reusable

## Test Utilities

### Mock Implementations
- `MockGuestbookService`: For controller testing
- `MockGuestbookRepository`: For service testing
- `MockValidator`: For custom validation testing

### Helper Functions
- `setupTestDB()`: Creates test database
- `setupController()`: Creates fully wired controller
- Factory methods for different configurations

## Running Specific Test Scenarios

### Test New Feature
```bash
go test -v -run TestNewFeature ./internal/service/...
```

### Test Error Handling
```bash
go test -v -run ".*Error.*" ./internal/...
```

### Test Validation
```bash
go test -v -run ".*Validation.*" ./internal/service/...
```

### Performance Testing
```bash
go test -v -bench=. ./internal/...
```

This comprehensive testing strategy ensures the refactored application maintains high quality while being easy to maintain and extend.
