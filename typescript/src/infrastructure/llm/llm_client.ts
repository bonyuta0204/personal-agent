import { LLMClient } from '../../usecase/generate_reply.ts';
import { DocumentRetriever } from './retriever_factory.ts';
import { ChatOpenAI } from '@langchain/openai';
import { StringOutputParser } from '@langchain/core/output_parsers';

export interface LLMConfig {
  apiKey?: string;
  model: string;
  temperature?: number;
  maxTokens?: number;
  retriever?: DocumentRetriever;
}

/**
 * Client for interacting with a language model API
 */
export class OpenAIClient implements LLMClient {
  private model: ChatOpenAI;
  private outputParser: StringOutputParser;
  
  constructor(private config: LLMConfig) {
    if (!config.apiKey && !Deno.env.get('OPENAI_API_KEY')) {
      throw new Error('OpenAI API key is required');
    }
    
    this.model = new ChatOpenAI({
      modelName: config.model,
      temperature: config.temperature || 0.7,
      maxTokens: config.maxTokens,
      openAIApiKey: config.apiKey || Deno.env.get('OPENAI_API_KEY'),
    });
    
    this.outputParser = new StringOutputParser();
  }

  async generateResponse(context: string | Record<string, unknown>): Promise<string> {
    try {
      // If we have a retriever, we would use it to augment the context with relevant documents
      if (this.config.retriever && typeof context === 'string') {
        const documents = await this.config.retriever.retrieveRelevantDocuments(context);
        // Augment context with retrieved documents
        // ...
      }
      
      // Convert context to proper format for LangChain
      const input = typeof context === 'string' 
        ? context 
        : JSON.stringify(context);
      
      // Call LangChain model
      const result = await this.model.invoke(input);
      return result.content.toString();
    } catch (error: unknown) {
      console.error('Error generating response:', error);
      const errorMessage = error instanceof Error ? error.message : String(error);
      throw new Error(`Failed to generate response: ${errorMessage}`);
    }
  }
}

/**
 * Factory for creating LLM clients
 */
export class LLMClientFactory {
  static createClient(type: 'openai', config: LLMConfig): LLMClient {
    switch (type) {
      case 'openai':
        return new OpenAIClient(config);
      default:
        throw new Error(`Unknown LLM client type: ${type}`);
    }
  }
}
