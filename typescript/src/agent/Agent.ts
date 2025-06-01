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
  createNewMemoryTool,
  retrieveMemoriesTool,
} from "../tools/memory.ts";
import { Pool } from "pg";

import { Config } from "../config/index.ts";

const systemMessage = new SystemMessage(
  [
    "You are an AI assistant that helps users find information from documents and manage personal memories. You have access to the following tools:",
    "",
    "## Document Search Tools",
    "",
    "You have access to several tools to help you find information from documents:",
    "",
    "- **Semantic Search (document_semantic_search)**: Finds relevant information based on the meaning of the query and document content. Useful for broad or open-ended questions, but may not always yield precise results.",
    "- **Tag Search (document_tag_search)**: Finds documents by matching tags. Use this when the user query matches known tags or when relevant tags are available—this often gives the best results for taggable topics.",
    "- **Keyword Search (document_keyword_search)**: Finds documents containing specific keywords. Use this for direct, simple queries or when other methods fail.",
    "",
    "## Memory Tools",
    "",
    "You can also manage personal memories:",
    "",
    "- **Create Memory (new_memory)**: Store a new memory with content, path, and tags for later retrieval.",
    "- **Retrieve Memories (retrieve_memories)**: Search and retrieve previously stored memories by path, tags, or get recent memories.",
    "",
    "**Tip:** Don’t rely solely on semantic search. Consider which tool best fits the user’s question and explain your reasoning. Tag search is particularly effective when relevant tags exist.",
    "",
    "## Response Format",
    "",
    "Respond directly to the user in a natural, conversational style. Summarize your reasoning and the steps you took to find the answer, mentioning which tools you used only if relevant. Cite sources or document details when appropriate. Avoid using a rigid structured format; your answer should read as a helpful, thoughtful reply to the user’s question.",
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
  const newMemoryTool = createNewMemoryTool(pool);
  const retrieveMemoriesToolInstance = retrieveMemoriesTool(pool);

  const agentTools = [
    documentSemanticTool,
    documentTagSearchTool,
    documentKeywordSearchTool,
    newMemoryTool,
    retrieveMemoriesToolInstance,
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
