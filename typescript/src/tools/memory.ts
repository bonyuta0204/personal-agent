import { z } from "zod";
import { DistanceStrategy, PGVectorStore } from "@langchain/community/vectorstores/pgvector";
import { OpenAIEmbeddings } from "@langchain/openai";
import { Pool, PoolConfig } from "pg";
import { tool } from "@langchain/core/tools";

import { Config } from "../config/index.ts";

export async function createMemorySemanticSearchTool(config: Config) {
  const vectorStoreConfig = {
    postgresConnectionOptions: {
      type: "postgres",
      host: config.database.host,
      port: config.database.port,
      user: config.database.username,
      password: config.database.password,
      database: config.database.database,
    } as PoolConfig,
    tableName: "memories",
    columns: {
      idColumnName: "id",
      vectorColumnName: "embedding",
      contentColumnName: "content",
    },
    distanceStrategy: "cosine" as DistanceStrategy,
  };

  const embeddings = new OpenAIEmbeddings({
    model: config.openai.embeddingModel,
  });

  const vectorStore = await PGVectorStore.initialize(
    embeddings,
    vectorStoreConfig,
  );

  return tool(
    async (input: { query: string; k?: number }) => {
      const retriever = vectorStore.asRetriever({
        k: input.k ?? 5,
        searchType: "similarity",
      });

      const results = await retriever.invoke(input.query);
      return JSON.stringify(results, null, 2);
    },
    {
      name: "memory_semantic_search",
      description: "Search agent memories by semantic similarity. Input: { query: string, k?: number }.",
      schema: z.object({ query: z.string(), k: z.number().optional() }),
    },
  );
}

export async function createMemoryTagSearchTool(pool: Pool) {
  // Fetch all unique tags from memories
  const client = await pool.connect();
  let tags: string[] = [];
  try {
    const res = await client.query(
      "SELECT DISTINCT jsonb_array_elements_text(tags) AS tag FROM memories WHERE jsonb_typeof(tags) = 'array';",
    );
    tags = res.rows.map((row: { tag: string }) => row.tag).sort();
  } finally {
    client.release();
  }
  const tagList = tags.length > 0 ? `\nExisting memory tags: [${tags.join(", ")}]` : "";

  return tool(
    async (input: { tags: string[]; k?: number }) => {
      const client = await pool.connect();
      try {
        const res = await client.query(
          `SELECT id, path, content, tags, created_at FROM memories WHERE tags @> $1::jsonb ORDER BY created_at DESC LIMIT $2`,
          [JSON.stringify(input.tags), input.k ?? 5],
        );
        return JSON.stringify(res.rows, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "memory_tag_search",
      description: `Search agent memories by tags. Input: { tags: string[], k?: number }.${tagList}`,
      schema: z.object({ tags: z.array(z.string()), k: z.number().optional() }),
    },
  );
}

export function createMemorySaveTool(pool: Pool, config: Config) {
  return tool(
    async (input: { path: string; content: string; tags?: string[] }) => {
      const client = await pool.connect();
      try {
        // Generate embedding using OpenAI
        const embeddings = new OpenAIEmbeddings({
          model: config.openai.embeddingModel,
        });
        const embedding = await embeddings.embedQuery(input.content);

        // Calculate SHA hash for content
        const crypto = await import("node:crypto");
        const sha = crypto.createHash("sha256").update(input.content).digest("hex");

        const tags = input.tags || [];
        const now = new Date();

        // Check if memory with same path already exists
        const existingRes = await client.query(
          "SELECT id, sha FROM memories WHERE path = $1",
          [input.path]
        );

        if (existingRes.rows.length > 0 && existingRes.rows[0].sha === sha) {
          return `Memory at path "${input.path}" already exists with same content.`;
        }

        let result;
        if (existingRes.rows.length > 0) {
          // Update existing memory
          result = await client.query(
            `UPDATE memories 
             SET content = $2, embedding = $3, tags = $4, sha = $5, modified_at = $6
             WHERE path = $1
             RETURNING id, path, created_at`,
            [input.path, input.content, JSON.stringify(embedding), JSON.stringify(tags), sha, now]
          );
        } else {
          // Create new memory
          result = await client.query(
            `INSERT INTO memories (path, content, embedding, tags, sha, modified_at)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING id, path, created_at`,
            [input.path, input.content, JSON.stringify(embedding), JSON.stringify(tags), sha, now]
          );
        }

        const memory = result.rows[0];
        return `Successfully saved memory: ID ${memory.id}, Path: ${memory.path}`;
      } finally {
        client.release();
      }
    },
    {
      name: "memory_save",
      description: "Save content as a memory. Input: { path: string, content: string, tags?: string[] }.",
      schema: z.object({ 
        path: z.string(), 
        content: z.string(), 
        tags: z.array(z.string()).optional() 
      }),
    },
  );
}

export function createMemoryListTool(pool: Pool) {
  return tool(
    async (input: { limit?: number; offset?: number }) => {
      const client = await pool.connect();
      try {
        const res = await client.query(
          `SELECT id, path, LEFT(content, 100) as content_preview, tags, created_at
           FROM memories 
           ORDER BY created_at DESC 
           LIMIT $1 OFFSET $2`,
          [input.limit ?? 10, input.offset ?? 0],
        );
        return JSON.stringify(res.rows, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "memory_list",
      description: "List recent memories. Input: { limit?: number, offset?: number }.",
      schema: z.object({ 
        limit: z.number().optional(), 
        offset: z.number().optional() 
      }),
    },
  );
}