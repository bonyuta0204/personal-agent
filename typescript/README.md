# Personal Agent – TypeScript Module

This module implements the AI agent for the Personal Agent project using Deno, LangChain, and PostgreSQL with pgvector. It provides a CLI interface for natural language chat and semantic document search.

## Features

- **AI Chat Agent**: Answers questions using retrieval-augmented generation (RAG) via LangChain
- **Semantic Document Search**: Uses OpenAI embeddings and pgvector for similarity search over documents
- **CLI Interface**: Interactive command-line chat with the agent
- **Configurable**: Reads DB and OpenAI settings from environment variables
- **Deno Runtime**: Secure, modern TypeScript runtime

## Architecture

```
src/
├── agent/      # AI agent logic (LangChain, tools)
├── cli/        # CLI entrypoint and loop
├── config/     # Loads config from env
├── tools/      # Document search tool (pgvector)
└── types/      # Type definitions
```

### Key Components

- **Agent**: LangChain-based agent with a semantic document search tool
- **CLI**: Text-based interface for chatting with the agent
- **Config**: Loads DB and OpenAI credentials from environment variables
- **Document Tool**: Connects to PostgreSQL/pgvector for semantic search

## Setup

1. **Environment Configuration**:
   
   Set the following environment variables (e.g., in a `.env` file or your shell):
   
   - `DB_HOST` (default: localhost)
   - `DB_PORT` (default: 5432)
   - `DB_NAME` (default: personal_agent)
   - `DB_USER` (default: postgres)
   - `DB_PASSWORD`
   - `DB_SSL` ("true" or "false")
   - `OPENAI_API_KEY` (required)
   - `OPENAI_MODEL` (default: gpt-4.1-mini)
   - `OPENAI_EMBEDDING_MODEL` (default: text-embedding-3-small)

2. **Database Setup**:
   - Ensure PostgreSQL is running with the `pgvector` extension enabled
   - Run database migrations (see Go module for schema)

3. **Dependencies**:
   - Managed via Deno import maps (`deno.json`)

## Usage

Start the interactive CLI:

```bash
deno task start
```

You’ll see:

```
🤖 Personal Agent
Type your questions naturally. Type 'exit' to quit.
```

Type any question or request. The agent will use semantic search over your documents and answer using retrieval-augmented generation.

**Example session:**

```
You: Find documents about vector search
🤔 ...

🤖 Agent:
I found 3 documents related to vector search. Here are the highlights:
- ...

You: exit
👋 Goodbye!
```

## Configuration

The agent loads configuration from environment variables. See `src/config/index.ts` for details. Example:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=personal_agent
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_SSL=false
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4.1-mini
OPENAI_EMBEDDING_MODEL=text-embedding-3-small
```

## How It Works

- The CLI launches an interactive loop (`src/cli/index.ts`)
- User input is sent to the agent (`src/agent/Agent.ts`)
- The agent uses LangChain and a semantic document search tool (`src/tools/document.ts`) to answer
- Results are retrieved from PostgreSQL with pgvector, using OpenAI embeddings
- The agent responds in natural language, optionally citing sources

## Development

```bash
deno task dev      # Watch mode

# Other tasks:
deno task test     # Run tests
deno task fmt      # Format code
deno task lint     # Lint code
deno task check    # Type check
```

## Integration with Go Module

This TypeScript module is designed to work alongside the Go CLI:
- **Shared Database**: Uses the same PostgreSQL schema as the Go collector/sync
- **Complementary**: Go handles data collection/sync, TypeScript provides chat & AI
