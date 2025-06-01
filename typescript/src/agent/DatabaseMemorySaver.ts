import { BaseCheckpointSaver } from "@langchain/langgraph";
import type { Checkpoint, CheckpointMetadata, CheckpointTuple } from "@langchain/langgraph";
import { ConversationRepository } from "../tools/conversation.ts";
import type { Client } from "https://deno.land/x/postgres@v0.19.3/mod.ts";

export class DatabaseMemorySaver extends BaseCheckpointSaver {
  private conversationRepo: ConversationRepository;

  constructor(private client: Client) {
    super();
    this.conversationRepo = new ConversationRepository(client);
  }

  async getTuple(config: { configurable?: { thread_id?: string; checkpoint_ns?: string; checkpoint_id?: string } }): Promise<CheckpointTuple | undefined> {
    const threadId = config.configurable?.thread_id;
    if (!threadId) return undefined;

    const messages = await this.conversationRepo.getConversationHistory(threadId);
    if (messages.length === 0) return undefined;

    // Convert database messages to checkpoint format
    const checkpoint: Checkpoint = {
      v: 1,
      id: threadId,
      ts: messages[messages.length - 1].created_at.toISOString(),
      channel_values: {},
      channel_versions: {},
      versions_seen: {},
      pending_sends: [],
    };

    // Store conversation history in channel_values
    const conversationHistory = messages.map(msg => ({
      role: msg.role,
      content: msg.content,
      metadata: msg.metadata,
      timestamp: msg.created_at.toISOString(),
    }));

    checkpoint.channel_values = {
      messages: conversationHistory,
    };

    const metadata: CheckpointMetadata = {
      source: "update" as const,
      step: messages.length,
      writes: {},
      parents: {},
    };

    return {
      config,
      checkpoint,
      metadata,
      parentConfig: undefined,
    };
  }

  async *list(
    config: { configurable?: { thread_id?: string; checkpoint_ns?: string } },
    options?: { before?: { configurable?: { checkpoint_id?: string } } }
  ): AsyncGenerator<CheckpointTuple> {
    const threadId = config.configurable?.thread_id;
    if (!threadId) return;

    const tuple = await this.getTuple(config);
    if (tuple) {
      yield tuple;
    }
  }

  async put(
    config: { configurable?: { thread_id?: string; checkpoint_ns?: string } },
    checkpoint: Checkpoint,
    metadata: CheckpointMetadata
  ): Promise<{ configurable: { thread_id?: string; checkpoint_id: string } }> {
    const threadId = config.configurable?.thread_id;
    if (!threadId) {
      throw new Error("thread_id is required in config.configurable");
    }

    // Extract thread info from thread_id format: "slack-{channel}-{user}" or "slack-{channel}-{user}-{thread_ts}"
    const parts = threadId.split('-');
    if (parts.length < 3 || parts[0] !== 'slack') {
      throw new Error("Invalid thread_id format. Expected: slack-{channel}-{user}[-{thread_ts}]");
    }

    const channelId = parts[1];
    const userId = parts[2];

    // Find or create conversation
    const conversation = await this.conversationRepo.findOrCreateConversation(
      threadId,
      channelId,
      userId
    );

    // Save new messages from checkpoint
    const messages = Array.isArray(checkpoint.channel_values?.messages) 
      ? checkpoint.channel_values.messages 
      : [];
    const existingMessages = await this.conversationRepo.getMessages(conversation.id);
    const existingCount = existingMessages.length;

    // Only save new messages (those after existing count)
    for (let i = existingCount; i < messages.length; i++) {
      const msg = messages[i];
      if (msg && typeof msg === 'object' && 'role' in msg && 'content' in msg) {
        await this.conversationRepo.addMessage(
          conversation.id,
          msg.role as 'human' | 'assistant' | 'system',
          msg.content as string,
          msg.metadata
        );
      }
    }

    return {
      configurable: {
        thread_id: threadId,
        checkpoint_id: checkpoint.id,
      },
    };
  }

  async putWrites(
    config: { configurable?: { thread_id?: string; checkpoint_ns?: string; checkpoint_id?: string } },
    writes: Array<[string, any]>,
    taskId: string
  ): Promise<void> {
    // For now, we'll handle writes through the put method
    // This could be extended to handle specific write operations
  }
}