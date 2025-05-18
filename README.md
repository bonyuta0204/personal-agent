# ✨ Personal Agent

📚 **Purpose**

> A personal agent that helps manage and interact with your documents and information sources.
> This repository contains tools for document synchronization, storage, and retrieval.

---

## 1. Tech Stack

| Layer                    | Runtime                                         | Libraries / Notes                               |
| ------------------------ | ----------------------------------------------- | ----------------------------------------------- |
| **Core Application**     | **Go**                                          | Clean architecture with domain-driven design    |
| **Document Management**  | **Go**                                          | Document synchronization and storage            |
| **Storage**              | GitHub, Local Storage                           | Multiple storage backends                       |
| **Deployment**           | Docker Compose (local)                          | Simple containerized deployment                 |

---

## 2. Repository Layout

```
personal-agent/
├─ go/                       # Go sources
│  ├─ internal/              # Private application code
│  │   ├─ domain/            # Enterprise business rules
│  │   │   ├─ model/         # Core domain entities and value objects
│  │   │   └─ port/          # Interfaces defining domain boundaries
│  │   ├─ usecase/           # Application business rules
│  │   │   ├─ document/      # Document-related use cases
│  │   │   └─ store/         # Storage-related use cases
│  │   └─ infrastructure/    # Frameworks, drivers, and external implementations
│  │       ├─ database/      # Database connections and utilities
│  │       ├─ embedding/     # Embedding service implementations
│  │       ├─ repository/    # Repository implementations
│  │       ├─ storage/       # Storage implementations (GitHub, local)
│  │       └─ util/          # Utility functions
│  ├─ cmd/                   # Application entry points
│  │   └─ cli/               # Command-line interface
│  ├─ config/                # Configuration files
│  ├─ migrations/            # Database migrations
│  └─ bin/                   # Compiled binaries
│
├─ docker-compose.yml        # Docker Compose configuration
└─ README.md (← **YOU ARE HERE**)
```

---

## 4. Quick Start (Local)

```bash
# 1. Clone the repository
git clone https://github.com/bonyuta0204/personal-agent.git
cd personal-agent

# 2. Setup environment variables
cp go/.env.sample go/.env

# 3. Build the CLI tool
cd go
make build

# 4. Run the CLI tool
./bin/cli --help
```

> **⚠️ Environment variables** are documented in `go/.env.sample`.

---

## 5. Storage Backends

The application supports multiple storage backends:

1. **GitHub Storage** - Store documents in a GitHub repository
2. **Local Storage** - Store documents locally on your machine

Storage implementations are located in `go/internal/infrastructure/storage/`.

---

## 6. Makefile Highlights

```makefile
build:             ## Build the CLI tool
test:              ## Run tests
clean:             ## Clean build artifacts
```

Check the `go/Makefile` for all available commands.

---

## 7. Document Management

The application provides commands for managing documents:

```bash
# List documents
./bin/cli document list

# Sync documents
./bin/cli document sync
```

Document operations are implemented in the `go/internal/usecase/document` package.

---

## 8. Project Architecture

The project follows clean architecture principles with a focus on domain-driven design:

1. **Domain Layer** - Contains core business entities and interfaces
2. **Use Case Layer** - Implements application-specific business rules
3. **Infrastructure Layer** - Provides concrete implementations of interfaces

This separation of concerns allows for easy testing and maintenance.

---

## 9. Contributing

1. Open a PR targeting `main`.
2. Ensure `make test lint` passes.
3. A reviewer merges after minimum one approval.

---

## 10. Roadmap

* [ ] Add support for additional storage backends
* [ ] Implement document versioning
* [ ] Add search functionality
* [ ] Improve CLI user experience

