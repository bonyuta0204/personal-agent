import { ChatSession, ChatSessionRepository } from '../domain/model/chat_session.ts';
import { Message, MessageRepository, MessageRole } from '../domain/model/message.ts';
import { ContextBuilder } from '../domain/service/context_builder.ts';

export interface GenerateReplyParams {
  sessionId: string;
}

export interface GenerateReplyResult {
  session: ChatSession;
  message: Message;
}

export interface LLMClient {
  generateResponse(context: string | Record<string, unknown>): Promise<string>;
}

export class GenerateReplyUseCase {
  constructor(
    private sessionRepo: ChatSessionRepository,
    private messageRepo: MessageRepository,
    private contextBuilder: ContextBuilder,
    private llmClient: LLMClient,
  ) {}

  async execute(params: GenerateReplyParams): Promise<GenerateReplyResult> {
    // Find the session
    const session = await this.sessionRepo.findById(params.sessionId);
    if (!session) {
      throw new Error(`Session not found: ${params.sessionId}`);
    }

    // Get all messages in the session
    const messages = await this.messageRepo.findBySessionId(session.id);
    
    // Build context for LLM
    const context = await this.contextBuilder.buildContext(session, messages);
    
    // Generate response using LLM
    const responseContent = await this.llmClient.generateResponse(context);
    
    // Create the assistant message
    const message = await this.messageRepo.create({
      sessionId: session.id,
      role: MessageRole.ASSISTANT,
      content: responseContent,
      createdAt: new Date(),
    });

    // Update the session's updatedAt timestamp
    session.updatedAt = new Date();
    await this.sessionRepo.update(session);

    return { session, message };
  }
}
