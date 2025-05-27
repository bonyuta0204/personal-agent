#!/usr/bin/env -S deno run --allow-all

import { Input } from "@cliffy/prompt";

import { DatabaseClient } from "../infra/pg/client.ts";
import { PersonalAgent, createPersonalAgent } from "../agent/Agent.ts";

class PersonalAgentCLI {
  private agent?: PersonalAgent;
  private db?: DatabaseClient;

  constructor() {
    // constructor(db: DatabaseClient) {
    this.agent = createPersonalAgent();
    // this.db = db;
  }

  async start(): Promise<void> {
    console.log("🤖 Personal Agent");
    console.log("Type your questions naturally. Type 'exit' to quit.\n");

    try {
      await this.conversationLoop();
    } catch (error) {
      console.error("❌ Failed to start:", error);
      Deno.exit(1);
    } finally {
      await this.cleanup();
    }
  }

  private async conversationLoop(): Promise<void> {
    while (true) {
      try {
        const input = await Input.prompt({
          message: "You:",
          minLength: 1,
        });

        const trimmedInput = input.trim();

        if (this.isExitCommand(trimmedInput)) {
          console.log("👋 Goodbye!");
          break;
        }

        await this.processInput(trimmedInput);
      } catch (error) {
        console.error("❌ Error:", error);
      }
    }
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
        { configurable: { thread_id: "personal_agent" } }
      );

      console.log("\n🤖 Agent:");
      console.debug(result);
      console.log(result.messages[result.messages.length - 1].content);
    } catch (error) {
      console.error("❌ Error:", error);
    }

    console.log(); // Add spacing
  }

  private async cleanup(): Promise<void> {
    if (this.db) {
      await this.db.disconnect();
    }
  }
}

// Main entry point
const cli = new PersonalAgentCLI();
await cli.start();
