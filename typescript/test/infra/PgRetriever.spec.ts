import { assertEquals, assertExists } from "https://deno.land/std@0.208.0/assert/mod.ts";
import { afterAll, beforeAll, describe, it } from "https://deno.land/std@0.208.0/testing/bdd.ts";

import { DatabaseClient } from "../../src/infra/pg/client.ts";
import { PgDocumentRetriever, PgMemoryRetriever } from "../../src/infra/pg/PgRetriever.ts";
import type { Document, Memory } from "../../src/types/database.ts";

describe("PgRetriever", () => {
  let db: DatabaseClient;
  let documentRetriever: PgDocumentRetriever;
  let memoryRetriever: PgMemoryRetriever;

  beforeAll(async () => {
    const config = {
      host: Deno.env.get("TEST_DB_HOST") || "localhost",
      port: parseInt(Deno.env.get("TEST_DB_PORT") || "5432"),
      database: Deno.env.get("TEST_DB_NAME") || "personal_agent_test",
      username: Deno.env.get("TEST_DB_USER") || "postgres",
      password: Deno.env.get("TEST_DB_PASSWORD") || "",
      ssl: false,
    };

    db = new DatabaseClient(config);
    await db.connect();

    documentRetriever = new PgDocumentRetriever(db);
    memoryRetriever = new PgMemoryRetriever(db);

    // Setup test data
    await setupTestData();
  });

  afterAll(async () => {
    // Cleanup test data
    await cleanupTestData();
    await db.disconnect();
  });

  async function setupTestData() {
    // Create a test store
    await db.query(
      "INSERT INTO stores (id, type, repo) VALUES (999, 'test', 'test-repo') ON CONFLICT (id) DO NOTHING",
    );

    // Create test documents
    const testDocuments: Partial<Document>[] = [
      {
        store_id: 999,
        path: "/test/doc1.md",
        content: "This is a test document about machine learning",
        embedding: Array(1536).fill(0.1),
        tags: ["ai", "ml"],
        sha: "test-sha-1",
      },
      {
        store_id: 999,
        path: "/test/doc2.md",
        content: "This document discusses web development",
        embedding: Array(1536).fill(0.2),
        tags: ["web", "frontend"],
        sha: "test-sha-2",
      },
    ];

    for (const doc of testDocuments) {
      await db.query(
        `INSERT INTO documents (store_id, path, content, embedding, tags, sha) 
         VALUES ($1, $2, $3, $4, $5, $6) 
         ON CONFLICT (path) DO NOTHING`,
        [
          doc.store_id,
          doc.path,
          doc.content,
          JSON.stringify(doc.embedding),
          JSON.stringify(doc.tags),
          doc.sha,
        ],
      );
    }

    // Create test memories
    const testMemories: Partial<Memory>[] = [
      {
        path: "/memory/mem1.md",
        content: "Important memory about project architecture",
        embedding: Array(1536).fill(0.3),
        tags: ["architecture", "important"],
        sha: "memory-sha-1",
      },
    ];

    for (const mem of testMemories) {
      await db.query(
        `INSERT INTO memories (path, content, embedding, tags, sha) 
         VALUES ($1, $2, $3, $4, $5) 
         ON CONFLICT (path) DO NOTHING`,
        [mem.path, mem.content, JSON.stringify(mem.embedding), JSON.stringify(mem.tags), mem.sha],
      );
    }
  }

  async function cleanupTestData() {
    await db.query("DELETE FROM documents WHERE store_id = 999");
    await db.query("DELETE FROM memories WHERE path LIKE '/memory/%'");
    await db.query("DELETE FROM stores WHERE id = 999");
  }

  describe("PgDocumentRetriever", () => {
    it("should find documents by path", async () => {
      const result = await documentRetriever.findByPath("/test/doc1.md");

      assertExists(result);
      assertEquals(result.path, "/test/doc1.md");
      assertEquals(result.store_id, 999);
    });

    it("should find documents by tags", async () => {
      const results = await documentRetriever.findByTags(["ai"]);

      assertEquals(results.length, 1);
      assertEquals(results[0].path, "/test/doc1.md");
    });

    it("should search documents by store ID", async () => {
      const results = await documentRetriever.searchByStoreId(
        999,
        "machine learning",
        { limit: 10, threshold: 0.0 },
      );

      assertEquals(results.length, 2);
      assertExists(results[0].similarity);
    });

    it("should search documents by embedding", async () => {
      const testEmbedding = Array(1536).fill(0.1);
      const results = await documentRetriever.searchByEmbedding(
        testEmbedding,
        { limit: 10, threshold: 0.0 },
      );

      assertEquals(results.length, 2);
      assertExists(results[0].similarity);
    });
  });

  describe("PgMemoryRetriever", () => {
    it("should find memories by path", async () => {
      const result = await memoryRetriever.findByPath("/memory/mem1.md");

      assertExists(result);
      assertEquals(result.path, "/memory/mem1.md");
    });

    it("should find memories by tags", async () => {
      const results = await memoryRetriever.findByTags(["architecture"]);

      assertEquals(results.length, 1);
      assertEquals(results[0].path, "/memory/mem1.md");
    });

    it("should search memories by embedding", async () => {
      const testEmbedding = Array(1536).fill(0.3);
      const results = await memoryRetriever.searchByEmbedding(
        testEmbedding,
        { limit: 10, threshold: 0.0 },
      );

      assertEquals(results.length, 1);
      assertExists(results[0].similarity);
    });
  });
});
