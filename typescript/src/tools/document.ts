import { z } from "zod";
import {
  DistanceStrategy,
  PGVectorStore,
} from "@langchain/community/vectorstores/pgvector";
import { OpenAIEmbeddings } from "@langchain/openai";
import { PoolConfig } from "pg";

import { Config } from "../config/index.ts";

export async function createDocumentSemanticTool(config: Config) {
  const vectorStoreConifg = {
    postgresConnectionOptions: {
      type: "postgres",
      host: config.database.host,
      port: config.database.port,
      user: config.database.username,
      password: config.database.password,
      database: config.database.database,
    } as PoolConfig,
    tableName: "documents",
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
    vectorStoreConifg
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
      name: "document_semantic_search",
      description:
        "Search for a document by its content. Input: { query: string, k?: number }.",
      schema: z.object({ query: z.string(), k: z.number().optional() }),
    }
  );
}

// タグによる検索ツール
import { Pool } from "pg";
import { tool } from "@langchain/core/tools";

export async function createDocumentTagSearchTool(pool: Pool) {
  // fetch all unique tags
  const client = await pool.connect();
  let tags: string[] = [];
  try {
    const res = await client.query(
      "SELECT DISTINCT jsonb_array_elements_text(tags) AS tag FROM documents WHERE jsonb_typeof(tags) = 'array';"
    );
    tags = res.rows.map((row: { tag: string }) => row.tag).sort();
  } finally {
    client.release();
  }
  const tagList =
    tags.length > 0 ? `\nExisting tags: [${tags.join(", ")}]` : "";
  return tool(
    async (input: { tags: string[]; k?: number }) => {
      const client = await pool.connect();
      try {
        const res = await client.query(
          `SELECT id, path, content, tags FROM documents WHERE tags @> $1::jsonb LIMIT $2`,
          [JSON.stringify(input.tags), input.k ?? 5]
        );
        return JSON.stringify(res.rows, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "document_tag_search",
      description: `Search documents by tags. Input: { tags: string[], k?: number }.${tagList}`,
      schema: z.object({ tags: z.array(z.string()), k: z.number().optional() }),
    }
  );
}

// キーワードによる検索ツール
export function createDocumentKeywordSearchTool(pool: Pool) {
  return tool(
    async (input: { keywords: string[]; k?: number }) => {
      const client = await pool.connect();
      try {
        const keywords = input.keywords.map((k) => `%${k}%`);
        const query = `
          SELECT DISTINCT id, path, content, tags 
          FROM documents 
          WHERE (
            ${Array(keywords.length)
              .fill("(content ILIKE $1 OR path ILIKE $1)")
              .join(" OR ")}
          )
          LIMIT $${keywords.length + 1}
        `;
        const res = await client.query(query, [...keywords, input.k ?? 5]);
        return JSON.stringify(res.rows, null, 2);
      } finally {
        client.release();
      }
    },
    {
      name: "document_keyword_search",
      description:
        "Search documents by keywords in content or path. Input: { keywords: string[], k?: number }.",
      schema: z.object({
        keywords: z.array(z.string()),
        k: z.number().optional(),
      }),
    }
  );
}
