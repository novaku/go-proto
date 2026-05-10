# SOLID Principles Implementation in Guestbook Application

This document explains how the SOLID principles have been implemented in the refactored guestbook application.

## Architecture Overview

The application has been restructured into multiple layers following clean architecture principles:

```
├── controller/          # Presentation Layer (gRPC handlers)
├── service/            # Business Logic Layer
├── repository/         # Data Access Layer
└── model/              # Domain Models
```

## SOLID Principles Implementation

### 1. Single Responsibility Principle (SRP)

**Each class/module has one reason to change.**

#### Examples:

**GuestbookController** (`internal/controller/guestbook_controller.go`)
- **Single Responsibility**: Handle gRPC request/response mapping
- **Only changes when**: gRPC interface needs modification

**GuestbookService** (`internal/service/guestbook_service.go`)
- **Single Responsibility**: Business logic for guestbook operations
- **Only changes when**: Business rules change

**GormGuestbookRepository** (`internal/repository/guestbook_repository.go`)
- **Single Responsibility**: Data persistence using GORM
- **Only changes when**: Database schema or GORM usage changes

**DefaultValidator** (`internal/service/validator.go`)
- **Single Responsibility**: Request validation
- **Only changes when**: Validation rules change

### 2. Open/Closed Principle (OCP)

**Software entities should be open for extension but closed for modification.**

#### Examples:

**Service Interface** (`internal/service/interfaces.go`)
```go
type GuestbookService interface {
    AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error)
    ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error)
}
```
- **Open for Extension**: Can create new implementations (e.g., CachedGuestbookService)
- **Closed for Modification**: Existing implementations don't need to change

**Repository Interface** (`internal/repository/interfaces.go`)
```go
type GuestbookRepository interface {
    Create(ctx context.Context, entry *model.GuestbookEntry) error
    FindWithPagination(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error)
}
```
- **Open for Extension**: Can add new data stores (Redis, MongoDB, etc.)
- **Closed for Modification**: Business logic doesn't change when data store changes

### 3. Liskov Substitution Principle (LSP)

**Objects of a superclass should be replaceable with objects of its subclasses.**

#### Examples:

**Repository Implementations**
- `GormGuestbookRepository` can be replaced with any other implementation
- All implementations must honor the contract defined by `GuestbookRepository`
- Factory pattern in `internal/controller/factory.go` allows easy substitution:

```go
func (f *GuestbookControllerFactory) CreateControllerWithCustomRepository(repo repository.GuestbookRepository) *GuestbookController {
    // Any implementation of GuestbookRepository can be used here
    validator := service.NewDefaultValidator()
    svc := service.NewGuestbookService(repo, validator)
    return NewGuestbookController(svc)
}
```

**Service Implementations**
- Different service implementations can be substituted without affecting the controller
- All must implement the `GuestbookService` interface correctly

### 4. Interface Segregation Principle (ISP)

**Clients should not be forced to depend on interfaces they do not use.**

#### Examples:

**Separated Interfaces**
- `GuestbookService`: Only guestbook-related operations
- `Validator`: Only validation operations
- `GuestbookRepository`: Only data access operations

**Small, Focused Interfaces**
```go
// Instead of one large interface, we have focused ones:

type GuestbookService interface {
    AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error)
    ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error)
}

type Validator interface {
    ValidateAddEntryRequest(req *pb.AddEntryRequest) error
    ValidateListEntriesRequest(req *pb.ListEntriesRequest) (*PaginationParams, error)
}
```

### 5. Dependency Inversion Principle (DIP)

**High-level modules should not depend on low-level modules. Both should depend on abstractions.**

#### Examples:

**Controller depends on Service abstraction**
```go
type GuestbookController struct {
    guestbookService service.GuestbookService // Interface, not concrete type
}
```

**Service depends on Repository abstraction**
```go
type guestbookService struct {
    repo      repository.GuestbookRepository // Interface, not concrete type
    validator Validator                      // Interface, not concrete type
}
```

**Dependency Injection in Factory**
```go
func (f *GuestbookControllerFactory) CreateController() *GuestbookController {
    // High-level factory creates and injects dependencies
    repo := repository.NewGormGuestbookRepository(f.db)
    validator := service.NewDefaultValidator()
    svc := service.NewGuestbookService(repo, validator)
    return NewGuestbookController(svc)
}
```

## Benefits of This Architecture

### 1. Testability
- Each layer can be tested in isolation using mocks
- Repository tests use in-memory SQLite
- Service tests use mock repositories
- Controller tests use mock services

### 2. Maintainability
- Changes in one layer don't affect others
- Easy to locate and fix bugs
- Clear separation of concerns

### 3. Extensibility
- New features can be added without modifying existing code
- New data stores can be added by implementing the repository interface
- New validation rules can be added by implementing the validator interface

### 4. Flexibility
- Different implementations can be swapped easily
- Configuration-driven behavior
- Support for different environments (test, dev, prod)

## Usage Examples

### Basic Usage
```go
// Create controller with default implementations
db := setupDatabase()
factory := controller.NewGuestbookControllerFactory(db)
ctrl := factory.CreateController()
```

### Custom Validator
```go
// Use custom validator
customValidator := &MyCustomValidator{}
ctrl := factory.CreateControllerWithCustomValidator(customValidator)
```

### Custom Repository
```go
// Use custom repository (e.g., Redis-based)
redisRepo := &RedisGuestbookRepository{client: redisClient}
ctrl := factory.CreateControllerWithCustomRepository(redisRepo)
```

## Testing Strategy

### Unit Tests
- Mock dependencies for fast, isolated tests
- Test business logic without database
- Test validation logic separately

### Integration Tests
- Test complete flows with real database
- Verify layer interactions
- Test error scenarios

### Test Coverage
- Controller: Request/response handling
- Service: Business logic and validation
- Repository: Data persistence and retrieval
- Validator: Input validation rules

This architecture ensures that the codebase is maintainable, testable, and extensible while following all SOLID principles.