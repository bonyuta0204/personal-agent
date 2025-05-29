import { ChatOpenAI } from "@langchain/openai";
import { MemorySaver } from "@langchain/langgraph";
import { createReactAgent } from "@langchain/langgraph/prebuilt";
import { SystemMessage } from "@langchain/core/messages";

import {
  createDocumentKeywordSearchTool,
  createDocumentSemanticTool,
  createDocumentTagSearchTool,
} from "../tools/document.ts";
import {
  createMemorySemanticSearchTool,
  createMemoryTagSearchTool,
  createMemorySaveTool,
  createMemoryListTool,
} from "../tools/memory.ts";
import { Pool } from "pg";

import { Config } from "../config/index.ts";

const systemMessage = new SystemMessage(
  [
    "You are an AI assistant that helps users find information from documents and manages personal memories. You have access to the following tools:",
    "",
    "## Document Search Tools",
    "",
    "You have access to several tools to help you find information from documents:",
    "",
    "- **Semantic Search (document_semantic_search)**: Finds relevant information based on the meaning of the query and document content. Useful for broad or open-ended questions, but may not always yield precise results.",
    "- **Tag Search (document_tag_search)**: Finds documents by matching tags. Use this when the user query matches known tags or when relevant tags are available—this often gives the best results for taggable topics.",
    "- **Keyword Search (document_keyword_search)**: Finds documents containing specific keywords. Use this for direct, simple queries or when other methods fail.",
    "",
    "## Memory Management Tools",
    "",
    "You also have access to memory management tools to store and retrieve personal information:",
    "",
    "- **Memory Semantic Search (memory_semantic_search)**: Search through saved memories using semantic similarity. Useful for finding related thoughts, notes, or experiences.",
    "- **Memory Tag Search (memory_tag_search)**: Find memories by specific tags. Effective when memories are categorized with tags.",
    "- **Memory Save (memory_save)**: Save new memories, thoughts, or important information. Use this when users want to remember something for later.",
    "- **Memory List (memory_list)**: List recent memories to see what has been saved previously.",
    "",
    "**Memory Usage:** When users mention wanting to remember something, ask to save it, or reference past conversations/thoughts, consider using memory tools. Always ask before saving sensitive or personal information.",
    "",
    "**Tip:** Don't rely solely on semantic search. Consider which tool best fits the user's question and explain your reasoning. Tag search is particularly effective when relevant tags exist.",
    "",
    "## Response Format",
    "",
    "Respond directly to the user in a natural, conversational style. Summarize your reasoning and the steps you took to find the answer, mentioning which tools you used only if relevant. Cite sources or document details when appropriate. Avoid using a rigid structured format; your answer should read as a helpful, thoughtful reply to the user's question.",
    "",
    "Begin!",
    "",
  ].join("\n"),
);

export type PersonalAgent = Awaited<ReturnType<typeof createPersonalAgent>>;

export async function createPersonalAgent(config: Config, pool: Pool) {
  // Define the tools for the agent to use
  const documentSemanticTool = await createDocumentSemanticTool(config);
  const documentTagSearchTool = await createDocumentTagSearchTool(pool);
  const documentKeywordSearchTool = createDocumentKeywordSearchTool(pool);
  
  // Memory tools
  const memorySemanticTool = await createMemorySemanticSearchTool(config);
  const memoryTagSearchTool = await createMemoryTagSearchTool(pool);
  const memorySaveTool = createMemorySaveTool(pool, config);
  const memoryListTool = createMemoryListTool(pool);

  const agentTools = [
    documentSemanticTool,
    documentTagSearchTool,
    documentKeywordSearchTool,
    memorySemanticTool,
    memoryTagSearchTool,
    memorySaveTool,
    memoryListTool,
  ];
  const agentModel = new ChatOpenAI({
    temperature: 0,
    model: config.openai.model,
  });

  // Initialize memory to persist state between graph runs
  const agentCheckpointer = new MemorySaver();
  const agent = createReactAgent({
    llm: agentModel,
    tools: agentTools,
    checkpointSaver: agentCheckpointer,
    prompt: systemMessage,
  });

  return agent;
}