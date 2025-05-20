---
trigger: glob
globs: go/**/*.go
---

# Go Architecture Guide

## Onion Architecture

We follow the onion architecture pattern in our Go codebase. This architecture is designed to create maintainable, testable, and loosely coupled applications by organizing code in concentric layers.

### Layers

Our implementation consists of the following layers (from innermost to outermost):

1. **Domain Layer** (`internal/domain/`)
   - Contains the business entities, value objects, and domain services
   - Defines interfaces that will be implemented by outer layers
   - Has no dependencies on other layers or external frameworks
   - Examples: Entity structs, repository interfaces, domain errors

2. **Use Case Layer** (`internal/usecase/`)
   - Implements application-specific business rules
   - Orchestrates the flow of data to and from entities
   - Depends only on the domain layer
   - Examples: Service implementations, business logic coordinators

3. **Adapter Layer**
   - Converts data between the format most convenient for use cases and entities
   - Implements interfaces defined in the domain layer
   - Examples: Repository implementations, API controllers, event handlers

4. **Infrastructure Layer** (`internal/infrastructure/`)
   - Contains frameworks, drivers, and tools like databases, web servers, etc.
   - Provides concrete implementations of the interfaces defined in inner layers
   - Examples: Database connections, external API clients, message queues

## Package Structure

```
go/
├── cmd/            # Application entry points
├── internal/
│   ├── domain/     # Domain entities and interfaces
│   ├── usecase/    # Application business logic
│   └── infrastructure/ # External tools and implementations
```

## Dependency Rule

The fundamental rule of the onion architecture is that dependencies always point inward. This means:

- Domain layer has no external dependencies
- Use case layer depends only on the domain layer
- Adapter layer depends on use case and domain layers
- Infrastructure layer depends on adapter, use case, and domain layers

## Best Practices

1. Use dependency injection to provide implementations of interfaces
2. Keep the domain layer clean and free of external dependencies
3. Use interfaces to define contracts between layers
4. Write unit tests for each layer independently
