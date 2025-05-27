# Personal Agent - TypeScript Module

TypeScript implementation of the Personal Agent using Deno, LangChain, and PostgreSQL with pgvector.

## Features

- **Vector Search**: Semantic search across documents and memories using OpenAI embeddings
- **LangChain Integration**: AI-powered query answering with retrieval-augmented generation (RAG)
- **PostgreSQL + pgvector**: High-performance vector database for embeddings
- **CLI Interface**: Command-line tool for querying and managing knowledge base
- **Deno Runtime**: Modern TypeScript runtime with built-in security and performance

## Architecture

### Directory Structure

```
src/
├── agent/           # AI agent implementation using LangChain
├── cli/             # Command-line interface
├── config/          # Configuration management
├── infra/pg/        # PostgreSQL infrastructure layer
├── tools/           # Retriever interfaces and tools
└── types/           # TypeScript type definitions

test/
└── infra/           # Infrastructure tests
```

### Key Components

- **Agent**: LangChain-powered AI agent for question answering
- **Retrievers**: Generic interfaces for vector search operations
- **PgRetriever**: PostgreSQL implementation with multiple search methods
- **DatabaseClient**: PostgreSQL client with connection management

## Setup

1. **Environment Configuration**:
   ```bash
   cp .env.example .env
   # Edit .env with your database and OpenAI credentials
   ```

2. **Database Setup**:
   - Ensure PostgreSQL is running with pgvector extension
   - Run migrations from the Go module to set up schema

3. **Dependencies**:
   All dependencies are managed through Deno's import maps in `deno.json`

## Usage

Simply run the command to start chatting:

```bash
deno task start
```

This starts an interactive chat session where you can:
- Ask questions naturally
- Search for documents and memories  
- Store new memories
- Get AI-powered answers with context

**Example Session:**
```bash
$ deno task start
🤖 Personal Agent
Type your questions naturally. Type 'exit' to quit.

✅ Ready to chat!

You: Find documents about machine learning
🤔 ...

🤖 Agent: I found 5 documents about machine learning. Here are the key insights:
[AI searches and provides comprehensive answer]

📚 Sources:
  📄 5 documents

You: Can you save a summary of the main points?
🤔 ...

🤖 Agent: I've saved a summary with the key machine learning concepts to memory.
[AI automatically stores the summary]

You: exit
👋 Goodbye!
```

### Natural Language Interface

The agent understands natural language and automatically:

1. **Analyzes Intent**: Understands what you want to do
2. **Selects Tools**: Chooses search, memory, or other tools
3. **Executes Actions**: Runs the appropriate operations  
4. **Provides Answers**: Gives comprehensive, contextual responses
5. **Learns**: Saves important interactions to memory

You can ask things like:
- "Find documents about X"
- "Remember that Y is important"
- "What did we discuss about Z?"
- "Search my memories for X"
- "Save this information: ..."

### Development Commands

```bash
# Run in development mode with file watching
deno task dev

# Run tests
deno task test

# Format code
deno task fmt

# Lint code
deno task lint

# Type check
deno task check
```

## Database Schema

The TypeScript module uses the same PostgreSQL schema as the Go module:

### documents

- Vector search with store filtering
- JSONB tags for flexible metadata
- SHA-based change detection

### memories

- Global vector search across all memories
- JSONB tags for categorization
- Independent of stores

## Integration with Go Module

This TypeScript module is designed to work alongside the Go CLI:

- **Shared Database**: Both modules use the same PostgreSQL schema
- **Compatible Types**: TypeScript types mirror Go domain models
- **Complementary Tools**: Go handles synchronization, TypeScript provides AI capabilities
