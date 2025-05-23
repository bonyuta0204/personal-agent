/**
 * ChatSession represents a conversation session with the personal agent
 */
export interface ChatSession {
  id: string;
  userId: string;
  title: string;
  createdAt: Date;
  updatedAt: Date;
  metadata?: Record<string, unknown>;
}

export interface ChatSessionRepository {
  create(session: Omit<ChatSession, "id">): Promise<ChatSession>;
  findById(id: string): Promise<ChatSession | null>;
  update(session: ChatSession): Promise<ChatSession>;
  delete(id: string): Promise<void>;
  listByUserId(userId: string): Promise<ChatSession[]>;
}
