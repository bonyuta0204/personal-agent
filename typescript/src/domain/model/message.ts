/**
 * Message represents a single message in a chat session
 */

export enum MessageRole {
  USER = 'user',
  ASSISTANT = 'assistant',
  SYSTEM = 'system'
}

export interface Message {
  id: string;
  sessionId: string;
  role: MessageRole;
  content: string;
  createdAt: Date;
  metadata?: Record<string, unknown>;
}

export interface MessageRepository {
  create(message: Omit<Message, 'id'>): Promise<Message>;
  findById(id: string): Promise<Message | null>;
  findBySessionId(sessionId: string): Promise<Message[]>;
  update(message: Message): Promise<Message>;
  delete(id: string): Promise<void>;
}
