#!/usr/bin/env -S deno run --allow-all

import { Input } from "@cliffy/prompt";

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
    console.log("Type your questions naturally. Type 'exit' to quit.\n");

    try {
      await this.conversationLoop();
    } catch (error) {
      console.error("❌ Failed to start:", error);
      Deno.exit(1);
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
        { configurable: { thread_id: "personal_agent" } },
      );

      console.log("\n🤖 Agent:");
      console.log(result.messages[result.messages.length - 1].content);
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
