import { Pool } from "pg";
import { tool } from "@langchain/core/tools";
import { z } from "zod";
import { OpenAIEmbeddings } from "@langchain/openai";
import { Config } from "../config/index.ts";

// Enhanced memory creation with automatic embedding generation
export function createNewMemoryTool(pool: Pool, config?: Config) {
  return tool(
    async (input: { 
      content: string; 
      path: string; 
      tags: string[];
      context?: string; // Additional context for better embedding
    }) => {
      const client = await pool.connect();
      try {
        let embedding = null;
        
        // Generate embedding if config is provided
        if (config) {
          const embeddings = new OpenAIEmbeddings({
            model: config.openai.embeddingModel,
          });
          const textToEmbed = input.context 
            ? `${input.context}\n\n${input.content}` 
            : input.content;
          embedding = await embeddings.embedQuery(textToEmbed);
        }
        
        const res = await client.query(
          embedding 
            ? `INSERT INTO memories (content, path, tags, embedding) VALUES ($1, $2, $3, $4::vector) RETURNING *`
            : `INSERT INTO memories (content, path, tags) VALUES ($1, $2, $3) RETURNING *`,
          embedding 
            ? [input.content, input.path, JSON.stringify(input.tags), `[${embedding.join(',')}]`]
            : [input.content, input.path, JSON.stringify(input.tags)]
        );
        
        return JSON.stringify({
          success: true,
          memory: res.rows[0],
          message: "Memory saved successfully"
        }, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "save_memory",
      description:
        "Save important information as a memory for future sessions. " +
        "Use this when the user shares personal preferences, important facts, " +
        "or information that should be remembered across conversations. " +
        "Input: { content: string, path: string, tags: string[], context?: string }",
      schema: z.object({
        content: z.string().describe("The main content to remember"),
        path: z.string().describe("Category path like 'preferences/coding' or 'facts/personal'"),
        tags: z.array(z.string()).describe("Relevant tags for easy retrieval"),
        context: z.string().optional().describe("Additional context to improve semantic search"),
      }),
    }
  );
}

export function retrieveMemoriesTool(pool: Pool) {
  return tool(
    async (input: { path?: string; tags?: string[]; limit?: number }) => {
      const client = await pool.connect();
      try {
        let query = "SELECT * FROM memories WHERE 1=1";
        const params: (string | number)[] = [];
        let paramCount = 0;

        if (input.path) {
          paramCount++;
          query += ` AND path = $${paramCount}`;
          params.push(input.path);
        }

        if (input.tags && input.tags.length > 0) {
          paramCount++;
          query += ` AND tags @> $${paramCount}::jsonb`;
          params.push(JSON.stringify(input.tags));
        }

        query += " ORDER BY created_at DESC";

        if (input.limit) {
          paramCount++;
          query += ` LIMIT $${paramCount}`;
          params.push(input.limit);
        }

        const res = await client.query(query, params);
        return JSON.stringify(res.rows, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "retrieve_memories",
      description:
        "Retrieve memories filtered by path and/or tags. Input: { path?: string, tags?: string[], limit?: number }.",
      schema: z.object({
        path: z.string().optional(),
        tags: z.array(z.string()).optional(),
        limit: z.number().optional(),
      }),
    }
  );
}

// Semantic memory search
export function createMemorySemanticSearchTool(pool: Pool, config: Config) {
  const embeddings = new OpenAIEmbeddings({
    model: config.openai.embeddingModel,
  });

  return tool(
    async (input: { query: string; k?: number; recencyWeight?: number }) => {
      const client = await pool.connect();
      try {
        // Generate embedding for the query
        const queryEmbedding = await embeddings.embedQuery(input.query);
        
        // Combine semantic similarity with recency
        const recencyWeight = input.recencyWeight || 0.1;
        const res = await client.query(
          `SELECT 
            id, path, content, tags, created_at,
            (1 - (embedding <=> $1::vector)) as similarity,
            ((1 - (embedding <=> $1::vector)) * (1 - $3) + 
             (1 / (1 + EXTRACT(EPOCH FROM (NOW() - created_at))/86400)) * $3) as combined_score
           FROM memories 
           WHERE embedding IS NOT NULL
           ORDER BY combined_score DESC 
           LIMIT $2`,
          [`[${queryEmbedding.join(',')}]`, input.k || 5, recencyWeight]
        );
        
        return JSON.stringify({
          query: input.query,
          memories: res.rows,
          count: res.rows.length
        }, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "search_memories_semantic",
      description: 
        "Search memories using semantic similarity. Useful for finding related " +
        "memories based on meaning rather than exact keywords. " +
        "Input: { query: string, k?: number, recencyWeight?: number (0-1) }",
      schema: z.object({
        query: z.string().describe("Natural language query to search memories"),
        k: z.number().optional().describe("Number of results to return (default: 5)"),
        recencyWeight: z.number().min(0).max(1).optional()
          .describe("Weight for recency vs similarity (0=only similarity, 1=only recency, default: 0.1)"),
      }),
    }
  );
}

// Update memory tool
export function createUpdateMemoryTool(pool: Pool) {
  return tool(
    async (input: { 
      id: number; 
      content?: string; 
      tags?: string[];
      appendContent?: boolean;
    }) => {
      const client = await pool.connect();
      try {
        let updateFields = [];
        let params = [];
        let paramCount = 0;

        if (input.content) {
          paramCount++;
          if (input.appendContent) {
            updateFields.push(`content = content || $${paramCount}`);
            params.push('\n\n' + input.content);
          } else {
            updateFields.push(`content = $${paramCount}`);
            params.push(input.content);
          }
        }

        if (input.tags) {
          paramCount++;
          updateFields.push(`tags = $${paramCount}::jsonb`);
          params.push(JSON.stringify(input.tags));
        }

        if (updateFields.length === 0) {
          return JSON.stringify({ error: "No fields to update" });
        }

        paramCount++;
        params.push(input.id);

        const res = await client.query(
          `UPDATE memories SET ${updateFields.join(', ')}, updated_at = NOW() 
           WHERE id = $${paramCount} RETURNING *`,
          params
        );

        return JSON.stringify({
          success: true,
          updated: res.rows[0]
        }, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "update_memory",
      description: 
        "Update an existing memory. Use this to add information to a memory " +
        "or correct/modify existing memories. " +
        "Input: { id: number, content?: string, tags?: string[], appendContent?: boolean }",
      schema: z.object({
        id: z.number().describe("Memory ID to update"),
        content: z.string().optional().describe("New content (replaces existing unless appendContent=true)"),
        tags: z.array(z.string()).optional().describe("New tags (replaces existing tags)"),
        appendContent: z.boolean().optional().describe("If true, appends content instead of replacing"),
      }),
    }
  );
}

// Memory analytics tool
export function createMemoryAnalyticsTool(pool: Pool) {
  return tool(
    async (input: { groupBy?: 'path' | 'tag' | 'date' }) => {
      const client = await pool.connect();
      try {
        let query: string;
        
        switch (input.groupBy) {
          case 'path':
            query = `
              WITH tag_expansion AS (
                SELECT path, jsonb_array_elements_text(tags) as tag
                FROM memories
              )
              SELECT 
                m.path, 
                COUNT(DISTINCT m.id) as count,
                array_agg(DISTINCT te.tag) as all_tags
              FROM memories m
              LEFT JOIN tag_expansion te ON m.path = te.path
              GROUP BY m.path
              ORDER BY count DESC`;
            break;
          case 'tag':
            query = `
              SELECT tag, COUNT(*) as count
              FROM memories, jsonb_array_elements_text(tags) as tag
              GROUP BY tag
              ORDER BY count DESC`;
            break;
          case 'date':
            query = `
              SELECT DATE(created_at) as date, COUNT(*) as count
              FROM memories
              GROUP BY DATE(created_at)
              ORDER BY date DESC
              LIMIT 30`;
            break;
          default:
            query = `
              SELECT 
                COUNT(*) as total_memories,
                COUNT(DISTINCT path) as unique_paths,
                (SELECT COUNT(DISTINCT tag) FROM memories, jsonb_array_elements_text(tags) as tag) as unique_tags,
                MIN(created_at) as oldest_memory,
                MAX(created_at) as newest_memory
              FROM memories`;
        }
        
        const res = await client.query(query);
        return JSON.stringify({
          groupBy: input.groupBy || 'summary',
          data: res.rows
        }, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "analyze_memories",
      description: 
        "Get analytics and insights about stored memories. " +
        "Useful for understanding memory patterns and organization. " +
        "Input: { groupBy?: 'path' | 'tag' | 'date' }",
      schema: z.object({
        groupBy: z.enum(['path', 'tag', 'date']).optional()
          .describe("Group memories by path, tag, or date. Omit for overall summary."),
      }),
    }
  );
}
