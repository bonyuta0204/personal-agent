import { ChatOpenAI } from "@langchain/openai";
import { MemorySaver } from "@langchain/langgraph";
import { createReactAgent } from "@langchain/langgraph/prebuilt";
import { RunnableToolLike } from "@langchain/core/runnables";

import { PromptTemplate } from "@langchain/core/prompts";

const promptTemplate = PromptTemplate.fromTemplate(
  // TODO: we have to brush up the prompt
  `You are a helpful assistant. Answer the user's question.
`
);

export type PersonalAgent = ReturnType<typeof createPersonalAgent>;

export function createPersonalAgent() {
  // Define the tools for the agent to use
  const agentTools: RunnableToolLike[] = [];
  const agentModel = new ChatOpenAI({ temperature: 0 });

  // Initialize memory to persist state between graph runs
  const agentCheckpointer = new MemorySaver();
  const agent = createReactAgent({
    llm: agentModel,
    tools: agentTools,
    checkpointSaver: agentCheckpointer,
    prompt: promptTemplate,
  });

  return agent;
}
