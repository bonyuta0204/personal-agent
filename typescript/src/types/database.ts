export type DocumentId = string;
export type MemoryId = string;
export type StoreId = number;

export interface Store {
  id: StoreId;
  type: string;
  repo?: string;
  created_at: Date;
  updated_at: Date;
}

export interface Document {
  id: DocumentId;
  store_id: StoreId;
  path: string;
  content: string;
  embedding: number[];
  tags: string[];
  sha: string;
  modified_at?: Date;
  created_at: Date;
  updated_at: Date;
}

export interface Memory {
  id: MemoryId;
  path: string;
  content: string;
  embedding: number[];
  tags: string[];
  sha: string;
  modified_at?: Date;
  created_at: Date;
  updated_at: Date;
}

export interface DocumentEntry {
  path: string;
  modified_at: Date;
}

export interface SearchResult<T> {
  item: T;
  similarity: number;
}

export interface SearchOptions {
  limit?: number;
  threshold?: number;
  tags?: string[] | undefined;
}
