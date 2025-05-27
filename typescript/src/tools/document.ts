import { tool } from "@langchain/core/tools";
import { z } from "zod";

import { DatabaseClient } from "../infra/pg/client.ts";

const documentSearchSchema = z.object({
  query: z.string().describe("The query to search for."),
});

type DocumentSearchInput = z.infer<typeof documentSearchSchema>;

export function createDocumentSemanticTool(db: DatabaseClient) {
  return tool(
    async ({ query }: DocumentSearchInput) => {
      const results = await db.queryRow<{ id: string; content: string }>(
        "SELECT id, content FROM documents WHERE content @@ to_tsquery($1) LIMIT 1",
        [query]
      );
      return results?.content || "No matching document found";
    },
    {
      name: "document_search",
      description: "Search for a document by its content.",
      schema: documentSearchSchema,
    }
  );
}
