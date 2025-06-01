import "jsr:@supabase/functions-js/edge-runtime.d.ts";
import { createHmac } from "https://deno.land/std@0.168.0/node/crypto.ts";
import { Pool } from "https://deno.land/x/postgres@v0.17.0/mod.ts";

// Import existing agent code directly
import { createPersonalAgent } from "../../../typescript/src/agent/Agent.ts";
import { Config } from "../../../typescript/src/config/index.ts";

// Slack signature verification
function verifySlackSignature(
  body: string,
  timestamp: string,
  signature: string,
  secret: string,
): boolean {
  const baseString = `v0:${timestamp}:${body}`;
  const hmac = createHmac("sha256", secret);
  hmac.update(baseString);
  const expectedSignature = `v0=${hmac.digest("hex")}`;
  return expectedSignature === signature;
}

// Format Slack message
function formatSlackResponse(agentResponse: string, sources?: any[]): any {
  const blocks = [
    {
      type: "section",
      text: {
        type: "mrkdwn",
        text: agentResponse,
      },
    },
  ];

  // Add sources if available
  if (sources && sources.length > 0) {
    blocks.push({
      type: "context",
      elements: sources.slice(0, 3).map((source) => ({
        type: "mrkdwn",
        text: `📄 ${source.path || source.type}`,
      })),
    });
  }

  return { blocks };
}

Deno.serve(async (req) => {
  try {
    // Get Slack headers
    const signature = req.headers.get("X-Slack-Signature");
    const timestamp = req.headers.get("X-Slack-Request-Timestamp");
    const body = await req.text();

    // Verify Slack signature
    const slackSecret = Deno.env.get("SLACK_SIGNING_SECRET");
    if (
      !slackSecret || !signature || !timestamp ||
      !verifySlackSignature(body, timestamp, signature, slackSecret)
    ) {
      return new Response("Unauthorized", { status: 401 });
    }

    // Parse Slack event
    const payload = JSON.parse(body);

    // Handle URL verification
    if (payload.type === "url_verification") {
      return new Response(payload.challenge);
    }

    // Handle events
    if (payload.type === "event_callback") {
      const event = payload.event;

      // Only respond to mentions and direct messages
      if (event.type !== "app_mention" && event.type !== "message") {
        return new Response("OK");
      }

      // Skip bot messages
      if (event.bot_id) {
        return new Response("OK");
      }

      // Extract message text (remove bot mention)
      const text = event.text.replace(/<@[A-Z0-9]+>/g, "").trim();

      // Create database connection (to existing PostgreSQL)
      const pool = new Pool({
        user: Deno.env.get("DB_USER"),
        password: Deno.env.get("DB_PASSWORD"),
        database: Deno.env.get("DB_NAME"),
        hostname: Deno.env.get("DB_HOST"),
        port: parseInt(Deno.env.get("DB_PORT") || "5432"),
        tls: { enabled: false }, // Adjust based on your setup
      }, 3);

      // Create config (reuse existing config structure)
      const config: Config = {
        openai: {
          apiKey: Deno.env.get("OPENAI_API_KEY")!,
          model: Deno.env.get("OPENAI_MODEL") || "gpt-4o-mini",
          embeddingModel: Deno.env.get("OPENAI_EMBEDDING_MODEL") ||
            "text-embedding-3-small",
        },
      };

      // Create agent instance
      const agent = await createPersonalAgent(config, pool);

      // Invoke agent
      const sessionId = `slack-${event.channel}-${event.user}`;
      const result = await agent.invoke(
        {
          messages: [{ role: "user", content: text }],
        },
        {
          configurable: { thread_id: sessionId },
        },
      );

      // Extract response
      const agentResponse = result.messages[result.messages.length - 1].content;

      // Post response to Slack
      const slackToken = Deno.env.get("SLACK_BOT_TOKEN");
      const slackResponse = await fetch(
        "https://slack.com/api/chat.postMessage",
        {
          method: "POST",
          headers: {
            "Authorization": `Bearer ${slackToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            channel: event.channel,
            thread_ts: event.thread_ts || event.ts,
            ...formatSlackResponse(agentResponse),
          }),
        },
      );

      if (!slackResponse.ok) {
        console.error("Failed to post to Slack:", await slackResponse.text());
      }

      // Clean up
      await pool.end();

      return new Response("OK");
    }

    return new Response("OK");
  } catch (error) {
    console.error("Error:", error);
    return new Response("Internal Server Error", { status: 500 });
  }
});
