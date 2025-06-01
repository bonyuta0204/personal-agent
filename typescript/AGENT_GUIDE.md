# Personal Agent Tool Interface & Prompt Guide

## Overview

Your personal agent is designed to effectively retrieve documents and memories, answer questions, and maintain context across sessions. This guide outlines the tool interfaces and best practices for optimal performance.

## Tool Interface Summary

### Document Search Tools

1. **document_semantic_search**
   - Purpose: Find conceptually related documents
   - Best for: Open-ended questions, broad topics
   - Input: `{ query: string, k?: number }`

2. **document_tag_search**
   - Purpose: Find documents by exact tag matches
   - Best for: Categorized content, specific topics
   - Input: `{ tags: string[], k?: number }`

3. **document_keyword_search**
   - Purpose: Find documents containing specific keywords
   - Best for: Exact terms, file paths, code snippets
   - Input: `{ keywords: string[], k?: number }`

### Memory Management Tools

1. **save_memory** (Enhanced)
   - Purpose: Store information with embeddings for future sessions
   - Features: Automatic embedding generation, context support
   - Input: `{ content: string, path: string, tags: string[], context?: string }`
   - Example paths: 'preferences/coding', 'projects/current', 'facts/personal'

2. **retrieve_memories**
   - Purpose: Get memories by path/tags filter
   - Best for: Structured retrieval
   - Input: `{ path?: string, tags?: string[], limit?: number }`

3. **search_memories_semantic** (New)
   - Purpose: Find memories using semantic similarity
   - Features: Recency weighting, similarity scoring
   - Input: `{ query: string, k?: number, recencyWeight?: number }`

4. **update_memory** (New)
   - Purpose: Modify existing memories
   - Features: Append or replace content
   - Input: `{ id: number, content?: string, tags?: string[], appendContent?: boolean }`

5. **analyze_memories** (New)
   - Purpose: Get insights about memory patterns
   - Features: Group by path, tag, or date
   - Input: `{ groupBy?: 'path' | 'tag' | 'date' }`

## Effective Usage Patterns

### 1. Conversation Initialization
```typescript
// At conversation start, the agent should:
1. Use analyze_memories() to get a summary
2. Use search_memories_semantic("user preferences") to personalize
3. Check for recent memories with retrieve_memories({ limit: 5 })
```

### 2. Information Retrieval Strategy
```typescript
// For user questions:
1. First: search_memories_semantic(query) - Check personal context
2. Then: document_tag_search() - Most precise if tags match
3. Fallback: document_semantic_search() - Broader search
4. Last resort: document_keyword_search() - Exact terms
```

### 3. Memory Creation Guidelines
```typescript
// Save memories for:
- User preferences: save_memory({
    content: "Prefers TypeScript over JavaScript",
    path: "preferences/coding",
    tags: ["typescript", "preferences"],
    context: "User mentioned while discussing project setup"
  })

- Project context: save_memory({
    content: "Working on personal AI agent with LangChain",
    path: "projects/current",
    tags: ["ai", "langchain", "current-project"]
  })

- Important facts: save_memory({
    content: "Uses PostgreSQL with pgvector for embeddings",
    path: "facts/technical",
    tags: ["database", "postgresql", "pgvector"]
  })
```

### 4. Memory Updates
```typescript
// Append new information:
update_memory({
  id: memoryId,
  content: "Also uses Deno for TypeScript runtime",
  appendContent: true
})

// Replace outdated info:
update_memory({
  id: memoryId,
  content: "Now prefers Bun over Deno",
  appendContent: false
})
```

## Prompt Engineering Tips

### System Prompt Structure

The enhanced system prompt includes:

1. **Role Definition**: Personal AI assistant with persistent memory
2. **Capability Overview**: Document search + memory management
3. **Memory Guidelines**: When/how to save and retrieve
4. **Search Strategy**: Multi-method approach with fallbacks
5. **Response Guidelines**: Natural, contextual, efficient

### Key Behavioral Instructions

1. **Proactive Memory Retrieval**
   - Always check memories at conversation start
   - Search memories before documents for personal context

2. **Intelligent Memory Creation**
   - Save preferences, project details, corrections
   - Use descriptive paths and relevant tags
   - Include context for better future retrieval

3. **Multi-Method Search**
   - Don't rely on single search method
   - Start specific (tags) then broaden (semantic)
   - Combine results from multiple sources

4. **Natural Responses**
   - Brief acknowledgment when saving memories
   - Cite sources when providing information
   - Maintain conversational tone

## Implementation Example

```typescript
// In your CLI or application:
import { createPersonalAgent } from "./agent/Agent_enhanced.ts";
import { Config } from "./config/index.ts";
import { Pool } from "pg";

// Initialize
const config = await loadConfig();
const pool = new Pool(config.database);
const agent = await createPersonalAgent(config, pool);

// Use the agent
const response = await agent.invoke({
  messages: [{ role: "user", content: "What's my preferred coding language?" }]
});

// The agent will:
// 1. Search memories for coding preferences
// 2. Provide personalized response based on stored memories
// 3. Update memories if new preferences are mentioned
```

## Best Practices

### Memory Path Organization
```
preferences/
  ├── coding/        # Language, framework preferences
  ├── tools/         # IDE, CLI tool preferences
  └── workflow/      # Development workflow preferences

projects/
  ├── current/       # Active projects
  ├── completed/     # Past projects
  └── ideas/         # Future project ideas

facts/
  ├── personal/      # User-specific information
  ├── technical/     # Technical details, configurations
  └── business/      # Business context, requirements

learning/
  ├── corrections/   # User corrections to agent responses
  ├── clarifications/# Clarified requirements
  └── examples/      # User-provided examples
```

### Tag Strategy
- Use hierarchical tags: ["coding", "typescript", "preferences"]
- Include temporal tags: ["current", "2024", "active"]
- Add context tags: ["project-x", "client-y", "personal"]

### Search Optimization
1. Build tag vocabulary from existing data
2. Use recency weight (0.1-0.3) for time-sensitive queries
3. Combine multiple search methods for comprehensive results
4. Cache frequently accessed memories

## Monitoring & Maintenance

Use `analyze_memories()` periodically to:
- Identify memory organization patterns
- Find duplicate or outdated memories
- Optimize tag usage
- Track memory growth over time

## Future Enhancements

Consider implementing:
1. Memory expiration/archival
2. Memory importance scoring
3. Cross-reference linking between memories
4. Memory templates for common patterns
5. Export/import memory snapshots