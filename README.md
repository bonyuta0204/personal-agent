# ✨ AI Project PM — Monorepo

📚 **Purpose**

> Centralise project knowledge (Slack, GitHub + Obsidian, Notion) and offer an "AI PM" that answers questions through a hybrid RAG pipeline.
> This repository contains *all* ingest, synchronisation, retrieval, and API layers.

---

## 1. Tech Stack

| Layer                    | Runtime                                         | Libraries / Notes                               |
| ------------------------ | ----------------------------------------------- | ----------------------------------------------- |
| **Ingest & Sync**        | **Go 1.22**                                     | `chi`, `pgx/v5`, `sqlc`, `pgmq-go`              |
| **RAG API & Embeddings** | **Node 20 (TypeScript)**                        | `LangChain 0.3`, `Fastify`, `@langchain/openai` |
| **Database**             | Postgres 16 + `pgvector 0.7` + `pgmq`           | HNSW index enabled                              |
| **CI/CD**                | GitHub Actions                                  | Go + Node matrix build                          |
| **Deployment**           | Docker Compose (local) · Helm charts (optional) |                                                 |

---

## 2. Repository Layout

```
repo-root/
├─ go/                       # Go sources
│  ├─ internal/              # Private application code
│  │   ├─ domain/            # Enterprise business rules
│  │   │   ├─ model/         # Core domain entities and value objects
│  │   │   └─ service/       # Domain services (pure business logic)
│  │   ├─ usecase/           # Application business rules
│  │   ├─ adapter/           # Interface adapters
│  │   └─ infrastructure/    # Frameworks, drivers, and external agency
│  ├─ pkg/                   # Public Go libraries
│  ├─ cmd/                   # Application entry points
│  │   ├─ ingest/            # Slack & GitHub (Obsidian) webhooks
│  │   ├─ sync/              # Notion ↔ memories delta sync
│  │   └─ migrate/           # DB migrations CLI
│  └─ db/                    # SQL migrations & sqlc.yaml
│
├─ node/                     # TypeScript sources
│  ├─ services/
│  │   ├─ rag-api/           # /query endpoint (Fastify + LangChain)
│  │   └─ embed-worker/      # Queue consumer → embeddings
│  └─ lib/                   # Shared TS utilities
│
├─ deploy/                   # docker‑compose, k8s manifests, Terraform
├─ .github/workflows/        # CI Pipelines
├─ Makefile                  # One‑liner dev commands
└─ README.md (← **YOU ARE HERE**)
```

---

## 3. Services

| Service          | Binary / Script              | Port   | Description                                                                        |
| ---------------- | ---------------------------- | ------ | ---------------------------------------------------------------------------------- |
| **ingest**       | `go/cmd/ingest`              | `8080` | Receives Slack & GitHub webhooks, stores raw events, enqueues embedding jobs.      |
| **sync**         | `go/cmd/sync`                | `8081` | Polls / Webhook from Notion, updates the `memories` table, emits embedding jobs.   |
| **migrate**      | `go/cmd/migrate`             | –      | Runs `golang-migrate`‐compatible SQL migrations.                                   |
| **embed‑worker** | `node/services/embed-worker` | –      | Dequeues from `pgmq`, calls embedding API, UPSERTs `documents`.                    |
| **rag‑api**      | `node/services/rag-api`      | `3000` | `/query` endpoint → LangChain retriever → LLM router (GPT‑4o / Claude 3 / Gemini). |

---

## 4. Quick Start (Local)

```bash
# 1. Clone & setup env vars
cp .env.example .env

# 2. Launch Postgres (with pgvector/pgmq) + services
make up           # docker‑compose up ‑d postgres

# 3. Apply database migrations
make migrate-up   # go run ./go/cmd/migrate up

# 4. Start Go services (live‑reload)
make run-ingest   # air -c .air.toml
make run-sync

# 5. Start Node services
cd node && pnpm dev        # runs rag‑api + embed‑worker via ts‑node‑dev
```

> **⚠️ Environment variables** are documented in `.env.example`.

---

## 5. Database Workflow

1. **Write migration** → `go/db/migrations/20240517_add_documents.up.sql`.
2. `make migrate-up` to apply locally (uses `golang-migrate`).
3. CI applies migrations against the test container; production uses the same binary.

SQL queries are generated into `go/internal/infrastructure/store` via **sqlc** for type‑safe access.

---

## 6. Makefile Highlights

```makefile
up:                ## Start postgres
migrate-up:        ## Apply latest migrations
run-ingest:        ## Run Go ingest with live reload
run-sync:          ## Run Go sync with live reload
run-node:          ## Start TS services (rag-api + worker)
lint test:         ## Static checks & unit tests
```

Run `make help` to list all targets.

---

## 7. CI Pipeline (GitHub Actions)

| Job                         | Purpose                                                             |
| --------------------------- | ------------------------------------------------------------------- |
| **go-test**                 | Lint (`golangci-lint`) & `go test ./...` against pgvector container |
| **node-test**               | `pnpm i` → ESLint + Vitest                                          |
| **docker‑build** (optional) | Build multi‑arch images for each service                            |

---

## 8. Deployment

1. **Docker Compose** for single‑host PoC.
2. **Kubernetes**: use `deploy/k8s/` manifests (Helm chart WIP).
3. **Supabase → RDS migration**: pg\_dump & restore; services rely only on `DATABASE_URL`.

---

## 9. Contributing

1. Open a PR targeting `main`.
2. Ensure `make test lint` passes.
3. A reviewer merges after minimum one approval.

---

## 10. Roadmap

* [ ] JWT‑based auth for `/query` (Slack‑signed): **Next**
* [ ] Add GitHub PR metadata to retrieval filters
* [ ] CI token‑usage dashboard
* [ ] Optional on‑prem Llama 3 inference

---

**Happy hacking 🚀**

