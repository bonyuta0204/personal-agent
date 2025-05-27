import { z } from "zod";
import {
  PGVectorStore,
  DistanceStrategy,
} from "@langchain/community/vectorstores/pgvector";
import { OpenAIEmbeddings } from "@langchain/openai";
import { PoolConfig } from "pg";

import { Config } from "../config/index.ts";

const documentSearchSchema = z.object({
  query: z.string().describe("The query to search for."),
});

type DocumentSearchInput = z.infer<typeof documentSearchSchema>;

export async function createDocumentSemanticTool(config: Config) {
  const vectorStoreConifg = {
    postgresConnectionOptions: {
      type: "postgres",
      host: config.database.host,
      port: config.database.port,
      user: config.database.username,
      password: config.database.password,
      database: config.database.database,
    } as PoolConfig,
    tableName: "documents",
    columns: {
      idColumnName: "id",
      vectorColumnName: "embedding",
      contentColumnName: "content",
    },
    distanceStrategy: "cosine" as DistanceStrategy,
  };

  const embeddings = new OpenAIEmbeddings({
    model: config.openai.embeddingModel,
  });

  const vectorStore = await PGVectorStore.initialize(
    embeddings,
    vectorStoreConifg
  );
  const retriever = vectorStore.asRetriever({
    k: 1,
    searchType: "similarity",
  });

  return retriever.asTool({
    name: "document_semantic_search",
    description: "Search for a document by its content.",
    schema: z.string(),
  });
}
