import { parse } from "https://deno.land/std@0.207.0/flags/mod.ts";
import { ChatSessionRepository } from "../domain/model/chat_session.ts";
import { MessageRepository } from "../domain/model/message.ts";
import { DefaultContextBuilder } from "../domain/service/context_builder.ts";
import { PgChatSessionRepository } from "../infrastructure/db/chat_session_repo_pg.ts";
import { PgMessageRepository } from "../infrastructure/db/message_repo_pg.ts";
import { LLMClientFactory } from "../infrastructure/llm/llm_client.ts";
import { StartSessionUseCase } from "../usecase/start_session.ts";
import { HandleUserMessageUseCase } from "../usecase/handle_user_message.ts";
import { GenerateReplyUseCase } from "../usecase/generate_reply.ts";

// Parse command line arguments
const args = parse(Deno.args, {
  string: ["session"],
  boolean: ["help"],
  alias: {
    s: "session",
    h: "help",
  },
});

// Show help
if (args.help) {
  console.log("Personal Agent CLI");
  console.log(
    "Usage: deno run --allow-env --allow-net src/cli/pm_chat.ts [options]"
  );
  console.log("Options:");
  console.log(
    "  -s, --session Session ID (optional, will create new session if not provided)"
  );
  console.log("  -h, --help    Show help");
  Deno.exit(0);
}

// Initialize dependencies
const dbClient = {}; // Replace with actual DB client initialization
const sessionRepo: ChatSessionRepository = new PgChatSessionRepository(
  dbClient
);
const messageRepo: MessageRepository = new PgMessageRepository(dbClient);
const contextBuilder = new DefaultContextBuilder();
const llmClient = LLMClientFactory.createClient("openai", {
  model: "gpt-4",
  temperature: 0.7,
});

// Initialize use cases
const startSessionUseCase = new StartSessionUseCase(sessionRepo, messageRepo);
const handleUserMessageUseCase = new HandleUserMessageUseCase(
  sessionRepo,
  messageRepo
);
const generateReplyUseCase = new GenerateReplyUseCase(
  sessionRepo,
  messageRepo,
  contextBuilder,
  llmClient
);

// Main function
async function main() {
  let sessionId = args.session;

  // Create a new session if none specified
  if (!sessionId) {
    console.log("Creating new chat session...");
    const session = await startSessionUseCase.execute({
      title: "CLI Chat Session",
      initialSystemMessage:
        "You are a helpful personal agent. Answer questions concisely and accurately.",
    });

    sessionId = session.id;
    console.log(`Session created with ID: ${sessionId}`);
  }

  console.log("Starting chat. Type 'exit' to quit.");

  // Chat loop
  while (true) {
    // Get user input
    const userInput = prompt("You: ");

    if (!userInput || userInput.toLowerCase() === "exit") {
      console.log("Goodbye!");
      break;
    }

    // Handle user message
    await handleUserMessageUseCase.execute({
      sessionId,
      content: userInput,
    });

    // Generate reply
    const { message } = await generateReplyUseCase.execute({ sessionId });

    console.log(`Assistant: ${message.content}`);
  }
}

// Run the main function
if (import.meta.main) {
  main().catch((error) => {
    console.error("Error:", error);
    Deno.exit(1);
  });
}
