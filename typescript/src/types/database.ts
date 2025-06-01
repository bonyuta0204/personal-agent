export type DocumentId = string;
export type MemoryId = string;
export type StoreId = number;
export type ConversationId = string;
export type MessageId = string;

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

export interface Conversation {
  id: ConversationId;
  thread_id: string;
  channel_id: string;
  user_id: string;
  created_at: Date;
  updated_at: Date;
}

export interface Message {
  id: MessageId;
  conversation_id: ConversationId;
  role: 'human' | 'assistant' | 'system';
  content: string;
  metadata?: Record<string, any>;
  created_at: Date;
}
