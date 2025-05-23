import {
  ChatSession,
  ChatSessionRepository as DomainChatSessionRepository,
} from "../domain/model/chat_session.ts";
import {
  Message,
  MessageRepository,
  MessageRole,
} from "../domain/model/message.ts";
// Remove unnecessary import
// import { ChatOpenAI } from "@langchain/openai";
// const model = new ChatOpenAI({ model: "gpt-4" });

export interface HandleUserMessageParams {
  sessionId: string;
  content: string;
  metadata?: Record<string, unknown>;
}

export interface HandleUserMessageResult {
  session: ChatSession;
  message: Message;
}

export class HandleUserMessageUseCase {
  constructor(
    private sessionRepo: DomainChatSessionRepository,
    private messageRepo: MessageRepository
  ) {}

  async execute(
    params: HandleUserMessageParams
  ): Promise<HandleUserMessageResult> {
    // Find the session
    const session = await this.sessionRepo.findById(params.sessionId);
    if (!session) {
      throw new Error(`Session not found: ${params.sessionId}`);
    }

    // Create the user message
    const message = await this.messageRepo.create({
      sessionId: session.id,
      role: MessageRole.USER,
      content: params.content,
      createdAt: new Date(),
      metadata: params.metadata,
    });

    // Update the session's updatedAt timestamp
    session.updatedAt = new Date();
    await this.sessionRepo.update(session);

    return { session, message };
  }
}
