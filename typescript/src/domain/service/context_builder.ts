import { ChatSession } from '../model/chat_session.ts';
import { Message } from '../model/message.ts';

/**
 * ContextBuilder is responsible for building the context for LLM interactions
 * based on the chat session and messages
 */
export interface ContextBuilder {
  /**
   * Builds a context for LLM interaction based on the chat session and messages
   * @param session The chat session
   * @param messages The messages in the session
   * @returns The context string or structured format required by the LLM
   */
  buildContext(session: ChatSession, messages: Message[]): Promise<string | Record<string, unknown>>;
}

/**
 * Default implementation of ContextBuilder
 */
export class DefaultContextBuilder implements ContextBuilder {
  async buildContext(session: ChatSession, messages: Message[]): Promise<string> {
    // Simple implementation that formats messages into a conversation format
    return messages
      .map((message) => `${message.role}: ${message.content}`)
      .join('\n\n');
  }
}
