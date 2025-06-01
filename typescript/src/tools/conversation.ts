import { Client } from "https://deno.land/x/postgres@v0.19.3/mod.ts";
import type { Conversation, ConversationId, Message, MessageId } from "../types/database.ts";

export class ConversationRepository {
  constructor(private client: Client) {}

  async findOrCreateConversation(
    threadId: string,
    channelId: string,
    userId: string
  ): Promise<Conversation> {
    const result = await this.client.queryObject<Conversation>`
      INSERT INTO conversations (thread_id, channel_id, user_id)
      VALUES (${threadId}, ${channelId}, ${userId})
      ON CONFLICT (thread_id) 
      DO UPDATE SET updated_at = CURRENT_TIMESTAMP
      RETURNING id, thread_id, channel_id, user_id, created_at, updated_at
    `;
    
    return result.rows[0];
  }

  async getConversation(threadId: string): Promise<Conversation | null> {
    const result = await this.client.queryObject<Conversation>`
      SELECT id, thread_id, channel_id, user_id, created_at, updated_at
      FROM conversations
      WHERE thread_id = ${threadId}
    `;
    
    return result.rows[0] || null;
  }

  async addMessage(
    conversationId: ConversationId,
    role: 'human' | 'assistant' | 'system',
    content: string,
    metadata?: Record<string, any>
  ): Promise<Message> {
    const result = await this.client.queryObject<Message>`
      INSERT INTO messages (conversation_id, role, content, metadata)
      VALUES (${conversationId}, ${role}, ${content}, ${metadata ? JSON.stringify(metadata) : null})
      RETURNING id, conversation_id, role, content, metadata, created_at
    `;
    
    return result.rows[0];
  }

  async getMessages(conversationId: ConversationId, limit = 100): Promise<Message[]> {
    const result = await this.client.queryObject<Message>`
      SELECT id, conversation_id, role, content, metadata, created_at
      FROM messages
      WHERE conversation_id = ${conversationId}
      ORDER BY created_at ASC
      LIMIT ${limit}
    `;
    
    return result.rows;
  }

  async getConversationHistory(threadId: string, limit = 100): Promise<Message[]> {
    const conversation = await this.getConversation(threadId);
    if (!conversation) {
      return [];
    }
    
    return this.getMessages(conversation.id, limit);
  }
}