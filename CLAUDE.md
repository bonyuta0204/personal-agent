# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Go Commands (from `/go` directory)
```bash
# Build
make build              # Build the CLI tool to bin/personal-agent
make all               # Run tests and build

# Testing & Quality
make test              # Run all tests with verbose output
make vet               # Run Go vet for suspicious constructs

# Development
make run               # Build and run the application
make clean             # Clean build artifacts

# Manual commands
go mod tidy            # Clean up dependencies
go build -o bin/personal-agent -v ./cmd/cli
```

### TypeScript Commands (from `/typescript` directory)
```bash
deno task start        # Run the AI chat CLI
deno task dev         # Run in watch mode
deno task test        # Run tests
deno task fmt         # Format code
deno task lint        # Lint code
deno task check       # Type check
```

### Docker Commands
```bash
docker-compose up -d   # Start PostgreSQL with pgvector
docker-compose down    # Stop services
```

## Architecture Overview

### Domain-Driven Design Structure

This project follows Clean Architecture with Domain-Driven Design (DDD) principles:

#### 1. Domain Layer (`internal/domain/`)
- **Models** (`model/`): Core business entities
  - `Document`: Represents knowledge base documents with embeddings
  - `Store`: Abstract document stores (currently GitHub)
  - `Memory`: Agent memory storage
- **Ports** (`port/`): Interface definitions that define domain boundaries
  - Repository interfaces for data persistence
  - Storage interfaces for external data sources
  - Embedding interfaces for AI/ML services

#### 2. Use Case Layer (`internal/usecase/`)
- **Document Operations** (`document/`): Document synchronization logic
- **Store Management** (`store/`): Store creation and management

#### 3. Infrastructure Layer (`internal/infrastructure/`)
- **Database** (`database/`): PostgreSQL connection management
- **Repositories** (`repository/postgres/`): Concrete repository implementations
- **Storage** (`storage/`): External storage implementations (GitHub)
- **Embedding** (`embedding/`): OpenAI embedding service integration
- **Utilities** (`util/`): Helper functions (e.g., SHA hashing)

### Key Architectural Patterns

#### 1. Hexagonal Architecture (Ports & Adapters)
- **Ports**: Abstract interfaces in `domain/port/`
- **Adapters**: Concrete implementations in `infrastructure/`
- Enables dependency inversion and testability

#### 2. Repository Pattern
- `DocumentRepository` and `StoreRepository` abstract data access
- PostgreSQL implementations in `infrastructure/repository/postgres/`

#### 3. Factory Pattern
- `StorageFactory` and `StorageFactoryProvider` for creating storage instances
- `EmbeddingProvider` factory for AI service abstraction

#### 4. Command Pattern (CLI)
- Cobra-based CLI with structured commands
- Commands: `store create/list`, `document sync`

### Database Schema

#### Core Tables (PostgreSQL + pgvector)
- **stores**: Document store configurations
- **documents**: Knowledge base documents with vector embeddings
- **memories**: Agent memory storage

#### Key Features
- **Vector Embeddings**: Uses pgvector extension for 1536-dimensional embeddings
- **JSONB Tags**: Flexible tagging system with GIN indexes
- **Automatic Timestamps**: Triggers for `updated_at` fields
- **SHA-based Change Detection**: Optimized sync by detecting unchanged documents

### Technology Stack

#### Backend (Go)
- **CLI Framework**: Cobra for command-line interface
- **Database**: PostgreSQL with pgvector extension
- **HTTP**: go-github for GitHub API integration
- **AI/ML**: OpenAI API for embeddings
- **Configuration**: godotenv for environment management

#### Frontend/AI (TypeScript/Deno)
- **Runtime**: Deno for secure TypeScript execution
- **AI Framework**: LangChain for agent orchestration
- **Chat Interface**: CLI-based interactive chat
- **Database Access**: Direct PostgreSQL connection for semantic search

### Key Conventions

#### 1. Error Handling
- Wrapped errors with context using `fmt.Errorf`
- Custom domain errors (e.g., `ErrUnsupportedStoreType`)
- Graceful failure with logging for non-critical operations

#### 2. Configuration Management
- Environment-based configuration with validation
- Default values for common settings (e.g., DB port 5432)
- Centralized config loading in `config/config.go`

#### 3. Data Flow
1. **Sync Process**: Storage → Document Fetching → SHA Comparison → Embedding Generation → Database Storage
2. **Change Detection**: Uses SHA hashing to avoid reprocessing unchanged documents
3. **Vector Embeddings**: Generated via OpenAI API for semantic search capabilities

#### 4. Code Organization
- Clear separation of concerns across layers
- Interface-driven design for testability
- Consistent naming conventions (e.g., `Repository`, `Provider`, `Factory`)

### Development Workflow

1. **Database Setup**: Use `docker-compose up -d` to start PostgreSQL with pgvector
2. **Environment Configuration**: Copy `.env.sample` to `.env` and configure credentials
3. **Build Go CLI**: Run `make build` in the `go/` directory
4. **Create Store**: `./bin/personal-agent store create owner/repo`
5. **Sync Documents**: `./bin/personal-agent document sync <store-id>`
6. **Run AI Agent**: Use `deno task start` in the `typescript/` directory for interactive chat

### Current Implementation Status

#### Implemented Features
- GitHub repository document synchronization
- Vector embedding generation and storage
- CLI interface for store and document management
- Change detection via SHA comparison
- PostgreSQL persistence with pgvector
- AI chat interface with RAG capabilities

#### Architecture Highlights
- **Extensible Storage**: Plugin-like architecture for different storage backends
- **AI-Ready**: Built-in support for vector embeddings and semantic search
- **Clean Boundaries**: Clear separation between domain logic and infrastructure
- **Type Safety**: Strong typing with custom types (e.g., `StoreId`, `DocumentId`)

The architecture demonstrates enterprise-level design patterns while maintaining simplicity and extensibility for future AI agent capabilities.