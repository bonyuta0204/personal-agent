export interface Config {
  database: {
    host: string;
    port: number;
    database: string;
    username: string;
    password: string;
    ssl: boolean;
  };
  openai: {
    openaiApiKey: string;
    model: string;
    embeddingModel: string;
  };
}

export function loadConfig(): Config {
  return {
    database: {
      host: Deno.env.get("DB_HOST") || "localhost",
      port: parseInt(Deno.env.get("DB_PORT") || "5432"),
      database: Deno.env.get("DB_NAME") || "personal_agent",
      username: Deno.env.get("DB_USER") || "postgres",
      password: Deno.env.get("DB_PASSWORD") || "",
      ssl: Deno.env.get("DB_SSL") === "true",
    },
    openai: {
      openaiApiKey: Deno.env.get("OPENAI_API_KEY") || "",
      model: Deno.env.get("OPENAI_MODEL") || "gpt-4.1-mini",
      embeddingModel: Deno.env.get("OPENAI_EMBEDDING_MODEL") || "text-embedding-3-small",
    },
  };
}
