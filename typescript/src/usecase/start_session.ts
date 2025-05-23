import { ChatSession, ChatSessionRepository } from '../domain/model/chat_session.ts';
import { Message, MessageRepository, MessageRole } from '../domain/model/message.ts';

export interface StartSessionParams {
  userId: string;
  title?: string;
  initialSystemMessage?: string;
}

export class StartSessionUseCase {
  constructor(
    private sessionRepo: ChatSessionRepository,
    private messageRepo: MessageRepository,
  ) {}

  async execute(params: StartSessionParams): Promise<ChatSession> {
    // Create a new chat session
    const session = await this.sessionRepo.create({
      userId: params.userId,
      title: params.title || 'New Conversation',
      createdAt: new Date(),
      updatedAt: new Date(),
    });

    // If an initial system message is provided, add it to the session
    if (params.initialSystemMessage) {
      await this.messageRepo.create({
        sessionId: session.id,
        role: MessageRole.SYSTEM,
        content: params.initialSystemMessage,
        createdAt: new Date(),
      });
    }

    return session;
  }
}
