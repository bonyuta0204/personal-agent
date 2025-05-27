import type { DatabaseClient } from "./client.ts";
import type {
  Document,
  Memory,
  SearchOptions,
  SearchResult,
} from "../../types/database.ts";
import type { DocumentRetriever, MemoryRetriever } from "../../tools/types.ts";

export class PgDocumentRetriever implements DocumentRetriever {
  constructor(private db: DatabaseClient) {}

  async search(
    query: string,
    options: SearchOptions = {}
  ): Promise<SearchResult<Document>[]> {
    const { limit = 10, threshold = 0.7 } = options;

    const sql = `
      SELECT d.*, 
             1 - (d.embedding <=> $1::vector) as similarity
      FROM documents d
      WHERE 1 - (d.embedding <=> $1::vector) > $2
      ${options.tags?.length ? "AND d.tags ?| $3" : ""}
      ORDER BY d.embedding <=> $1::vector
      LIMIT $${options.tags?.length ? "4" : "3"}
    `;

    const params: unknown[] = [query, threshold];
    if (options.tags?.length) {
      params.push(options.tags);
    }
    params.push(limit);

    const rows = await this.db.query<Document & { similarity: number }>(
      sql,
      params
    );

    return rows.map((row) => ({
      item: {
        id: row.id,
        store_id: row.store_id,
        path: row.path,
        content: row.content,
        embedding: row.embedding,
        tags: row.tags,
        sha: row.sha,
        modified_at: row.modified_at || new Date(),
        created_at: row.created_at,
        updated_at: row.updated_at,
      },
      similarity: row.similarity,
    }));
  }

  async searchByEmbedding(
    embedding: number[],
    options: SearchOptions = {}
  ): Promise<SearchResult<Document>[]> {
    const { limit = 10, threshold = 0.7 } = options;

    const sql = `
      SELECT d.*, 
             1 - (d.embedding <=> $1::vector) as similarity
      FROM documents d
      WHERE 1 - (d.embedding <=> $1::vector) > $2
      ${options.tags?.length ? "AND d.tags ?| $3" : ""}
      ORDER BY d.embedding <=> $1::vector
      LIMIT $${options.tags?.length ? "4" : "3"}
    `;

    const params: unknown[] = [JSON.stringify(embedding), threshold];
    if (options.tags?.length) {
      params.push(options.tags);
    }
    params.push(limit);

    const rows = await this.db.query<Document & { similarity: number }>(
      sql,
      params
    );

    return rows.map((row) => ({
      item: {
        id: row.id,
        store_id: row.store_id,
        path: row.path,
        content: row.content,
        embedding: row.embedding,
        tags: row.tags,
        sha: row.sha,
        modified_at: row.modified_at || new Date(),
        created_at: row.created_at,
        updated_at: row.updated_at,
      },
      similarity: row.similarity,
    }));
  }

  async searchByStoreId(
    storeId: number,
    query: string,
    options: SearchOptions = {}
  ): Promise<SearchResult<Document>[]> {
    const { limit = 10, threshold = 0.7 } = options;

    const sql = `
      SELECT d.*, 
             1 - (d.embedding <=> $1::vector) as similarity
      FROM documents d
      WHERE d.store_id = $2
        AND 1 - (d.embedding <=> $1::vector) > $3
      ${options.tags?.length ? "AND d.tags ?| $4" : ""}
      ORDER BY d.embedding <=> $1::vector
      LIMIT $${options.tags?.length ? "5" : "4"}
    `;

    const params: unknown[] = [query, storeId, threshold];
    if (options.tags?.length) {
      params.push(options.tags);
    }
    params.push(limit);

    const rows = await this.db.query<Document & { similarity: number }>(
      sql,
      params
    );

    return rows.map((row) => ({
      item: {
        id: row.id,
        store_id: row.store_id,
        path: row.path,
        content: row.content,
        embedding: row.embedding,
        tags: row.tags,
        sha: row.sha,
        modified_at: row.modified_at || new Date(),
        created_at: row.created_at,
        updated_at: row.updated_at,
      },
      similarity: row.similarity,
    }));
  }

  async findByPath(path: string): Promise<Document | null> {
    const sql = "SELECT * FROM documents WHERE path = $1";
    const result = await this.db.queryRow<Document>(sql, [path]);
    if (result && result.modified_at === undefined) {
      result.modified_at = new Date();
    }
    return result;
  }

  async findByTags(
    tags: string[],
    options: SearchOptions = {}
  ): Promise<Document[]> {
    const { limit = 10 } = options;

    const sql = `
      SELECT * FROM documents 
      WHERE tags ?| $1
      ORDER BY created_at DESC
      LIMIT $2
    `;

    const results = await this.db.query<Document>(sql, [tags, limit]);
    return results.map((result) => ({
      ...result,
      modified_at: result.modified_at || new Date(),
    }));
  }
}

export class PgMemoryRetriever implements MemoryRetriever {
  constructor(private db: DatabaseClient) {}

  async search(
    query: string,
    options: SearchOptions = {}
  ): Promise<SearchResult<Memory>[]> {
    const { limit = 10, threshold = 0.7 } = options;

    const sql = `
      SELECT m.*, 
             1 - (m.embedding <=> $1::vector) as similarity
      FROM memories m
      WHERE 1 - (m.embedding <=> $1::vector) > $2
      ${options.tags?.length ? "AND m.tags ?| $3" : ""}
      ORDER BY m.embedding <=> $1::vector
      LIMIT $${options.tags?.length ? "4" : "3"}
    `;

    const params: unknown[] = [query, threshold];
    if (options.tags?.length) {
      params.push(options.tags);
    }
    params.push(limit);

    const rows = await this.db.query<Memory & { similarity: number }>(
      sql,
      params
    );

    return rows.map((row) => ({
      item: {
        id: row.id,
        path: row.path,
        content: row.content,
        embedding: row.embedding,
        tags: row.tags,
        sha: row.sha,
        modified_at: row.modified_at || new Date(),
        created_at: row.created_at,
        updated_at: row.updated_at,
      },
      similarity: row.similarity,
    }));
  }

  async searchByEmbedding(
    embedding: number[],
    options: SearchOptions = {}
  ): Promise<SearchResult<Memory>[]> {
    const { limit = 10, threshold = 0.7 } = options;

    const sql = `
      SELECT m.*, 
             1 - (m.embedding <=> $1::vector) as similarity
      FROM memories m
      WHERE 1 - (m.embedding <=> $1::vector) > $2
      ${options.tags?.length ? "AND m.tags ?| $3" : ""}
      ORDER BY m.embedding <=> $1::vector
      LIMIT $${options.tags?.length ? "4" : "3"}
    `;

    const params: unknown[] = [JSON.stringify(embedding), threshold];
    if (options.tags?.length) {
      params.push(options.tags);
    }
    params.push(limit);

    const rows = await this.db.query<Memory & { similarity: number }>(
      sql,
      params
    );

    return rows.map((row) => ({
      item: {
        id: row.id,
        path: row.path,
        content: row.content,
        embedding: row.embedding,
        tags: row.tags,
        sha: row.sha,
        modified_at: row.modified_at || new Date(),
        created_at: row.created_at,
        updated_at: row.updated_at,
      },
      similarity: row.similarity,
    }));
  }

  async findByPath(path: string): Promise<Memory | null> {
    const sql = "SELECT * FROM memories WHERE path = $1";
    const result = await this.db.queryRow<Memory>(sql, [path]);
    if (result && result.modified_at === undefined) {
      result.modified_at = new Date();
    }
    return result;
  }

  async findByTags(
    tags: string[],
    options: SearchOptions = {}
  ): Promise<Memory[]> {
    const { limit = 10 } = options;

    const sql = `
      SELECT * FROM memories 
      WHERE tags ?| $1
      ORDER BY created_at DESC
      LIMIT $2
    `;

    const results = await this.db.query<Memory>(sql, [tags, limit]);
    return results.map((result) => ({
      ...result,
      modified_at: result.modified_at || new Date(),
    }));
  }
}
