import { Pool } from "pg";
import { tool } from "@langchain/core/tools";
import { z } from "zod";

export function createNewMemoryTool(pool: Pool) {
  return tool(
    async (input: { content: string; path: string; tags: string[] }) => {
      const client = await pool.connect();
      try {
        const res = await client.query(
          `INSERT INTO memories (content, path, tags) VALUES ($1, $2, $3) RETURNING *`,
          [input.content, input.path, JSON.stringify(input.tags)]
        );
        return JSON.stringify(res.rows, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "new_memory",
      description:
        "Create a new memory. Input: { content: string, path: string, tags: string[] }.",
      schema: z.object({
        content: z.string(),
        path: z.string(),
        tags: z.array(z.string()),
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
