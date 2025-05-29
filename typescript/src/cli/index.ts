#!/usr/bin/env -S deno run --allow-all

import { readLines } from "https://deno.land/std@0.208.0/io/mod.ts";

import { createPersonalAgent, PersonalAgent } from "../agent/Agent.ts";
import { loadConfig } from "../config/index.ts";
import { Pool } from "pg";

class PersonalAgentCLI {
  private agent: PersonalAgent;

  constructor(agent: PersonalAgent) {
    this.agent = agent;
  }

  async start(): Promise<void> {
    console.log("🤖 Personal Agent");
    console.log("Type your questions naturally.");
    console.log("Exit: Type 'exit' or '終了', or press Ctrl+C\n");

    try {
      await this.conversationLoop();
    } catch (error) {
      console.error("❌ Failed to start:", error);
      Deno.exit(1);
    }
  }

  private async conversationLoop(): Promise<void> {
    // Set up signal handlers for graceful shutdown
    const handleInterrupt = () => {
      console.log("\n\n👋 Goodbye!");
      Deno.exit(0);
    };

    Deno.addSignalListener("SIGINT", handleInterrupt);
    Deno.addSignalListener("SIGTERM", handleInterrupt);

    const encoder = new TextEncoder();
    const lineReader = readLines(Deno.stdin);

    while (true) {
      try {
        // Display prompt
        await Deno.stdout.write(encoder.encode("You: "));

        // Read user input
        const result = await lineReader.next();
        if (result.done) {
          console.log("\n👋 Goodbye!");
          break;
        }

        const input = result.value;
        const trimmedInput = input.trim();

        if (this.isExitCommand(trimmedInput)) {
          console.log("👋 Goodbye!");
          break;
        }

        if (trimmedInput.length > 0) {
          await this.processInput(trimmedInput);
        }
      } catch (error) {
        // Check if the error is due to EOF/Ctrl+D
        if (
          error instanceof Error &&
          (error.message.includes("EOF") || error.message.includes("Bad resource ID"))
        ) {
          console.log("\n👋 Goodbye!");
          break;
        }
        console.error("❌ Error:", error);
      }
    }

    // Clean up signal listeners
    Deno.removeSignalListener("SIGINT", handleInterrupt);
    Deno.removeSignalListener("SIGTERM", handleInterrupt);
  }

  private isExitCommand(input: string): boolean {
    const exitCommands = ["exit", "quit", "bye", "goodbye", "終了"];
    return exitCommands.includes(input.toLowerCase());
  }

  private async processInput(input: string): Promise<void> {
    if (!this.agent) return;

    console.log("\n🤔 ...");

    try {
      const result = await this.agent.invoke(
        {
          messages: [
            {
              role: "user",
              content: input,
            },
          ],
        },
        { configurable: { thread_id: "personal_agent" } },
      );

      console.log("\n🤖 Agent:");
      const lastMessage = result.messages[result.messages.length - 1];
      const content = lastMessage.content;

      // Ensure proper UTF-8 output
      if (typeof content === "string") {
        // Force UTF-8 encoding for console output
        const encoder = new TextEncoder();
        const decoder = new TextDecoder("utf-8");
        const encoded = encoder.encode(content);
        const decoded = decoder.decode(encoded);
        console.log(decoded);
      } else {
        console.log(content);
      }
    } catch (error) {
      console.error("❌ Error:", error);
    }

    console.log(); // Add spacing
  }
}

// Main entry point
const config = await loadConfig();
const pool = new Pool({
  host: config.database.host,
  port: config.database.port,
  user: config.database.username,
  password: config.database.password,
  database: config.database.database,
  ssl: config.database.ssl,
});
const personalAgent = await createPersonalAgent(config, pool);
const cli = new PersonalAgentCLI(personalAgent);
await cli.start();
