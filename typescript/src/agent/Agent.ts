import { ChatOpenAI } from "@langchain/openai";
import { MemorySaver } from "@langchain/langgraph";
import { createReactAgent } from "@langchain/langgraph/prebuilt";
import { SystemMessage } from "@langchain/core/messages";

import { createDocumentSemanticTool } from "../tools/document.ts";

import { Config } from "../config/index.ts";

const systemMessage = new SystemMessage(
  `
You are an AI assistant that helps users find information from documents. You have access to the following tools:


## Document Search Tool
- Use the 'document_semantic_search' tool to find relevant information from documents when:
  - The user asks a question that might be answered by stored documents
  - You need to verify information before responding
  - The user asks for specific details that might be in documents
- The tool performs semantic search, so use natural language queries
- Always analyze the document content carefully before responding


## Response Format
Use the following format:

Question: the input question you must answer
Thought: you should always think about what to do
Action: the action to take, should be one of 
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question, citing sources when appropriate

Begin!

`
);

export type PersonalAgent = Awaited<ReturnType<typeof createPersonalAgent>>;

export async function createPersonalAgent(config: Config) {
  // Define the tools for the agent to use
  const documentSemanticTool = await createDocumentSemanticTool(config);

  const agentTools = [documentSemanticTool];
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
