import "jsr:@supabase/functions-js/edge-runtime.d.ts";
import { createHmac } from "https://deno.land/std@0.168.0/node/crypto.ts";
import { Pool } from "pg";

import { createPersonalAgent } from "../../../typescript/src/agent/Agent.ts";
import { Config } from "../../../typescript/src/config/index.ts";

// Import EdgeRuntime for background tasks
declare const EdgeRuntime: {
  waitUntil(promise: Promise<any>): void;
};

// Slack signature verification
function verifySlackSignature(
  body: string,
  timestamp: string,
  signature: string,
  secret: string,
): boolean {
  const baseString = `v0:${timestamp}:${body}`;
  const hmac = createHmac("sha256", secret);
  console.log("[baseString]", baseString);
  console.log("[signature]", signature);
  console.log("[secret]", secret);
  hmac.update(baseString);
  const expectedSignature = `v0=${hmac.digest("hex")}`;
  console.log("[expectedSignature]", expectedSignature);
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

// Process Slack message in background
async function processSlackMessage(event: any, text: string) {
  let pool: Pool | null = null;
  
  try {
    console.log(`Processing message from ${event.user} in channel ${event.channel}: ${text}`);
    
    // Create database connection
    pool = new Pool({
      connectionString: Deno.env.get("SUPABASE_DB_URL")!,
    });

    // Create config
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
    } else {
      console.log("Successfully posted response to Slack");
    }
  } catch (error) {
    console.error("Error processing Slack message:", error);
    
    // Try to send error message to Slack
    try {
      const slackToken = Deno.env.get("SLACK_BOT_TOKEN");
      await fetch(
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
            text: "Sorry, I encountered an error while processing your message. Please try again later.",
          }),
        },
      );
    } catch (postError) {
      console.error("Failed to post error message to Slack:", postError);
    }
  } finally {
    // Clean up database connection
    if (pool) {
      await pool.end();
    }
  }
}

Deno.serve(async (req) => {
  console.log("Request received:", req.method, req.url);
  console.log("Headers:", Object.fromEntries(req.headers.entries()));
  
  try {
    // Get Slack headers
    const signature = req.headers.get("X-Slack-Signature");
    const timestamp = req.headers.get("X-Slack-Request-Timestamp");
    const contentType = req.headers.get("Content-Type");
    const body = await req.text();

    // Log environment variables (without exposing secrets)
    console.log("Environment check:");
    console.log("- SLACK_SIGNING_SECRET exists:", !!Deno.env.get("SLACK_SIGNING_SECRET"));
    console.log("- SLACK_BOT_TOKEN exists:", !!Deno.env.get("SLACK_BOT_TOKEN"));
    
    // Verify Slack signature
    const slackSecret = Deno.env.get("SLACK_SIGNING_SECRET");
    console.log("Verification details:");
    console.log("- Signature header:", signature);
    console.log("- Timestamp header:", timestamp);
    console.log("- Content-Type:", contentType);
    console.log("- Body length:", body.length);
    console.log("- Body preview:", body.substring(0, 100));
    
    if (!slackSecret) {
      console.error("SLACK_SIGNING_SECRET is not set!");
      return new Response("Server configuration error", { status: 500 });
    }
    
    if (!signature || !timestamp) {
      console.error("Missing required headers");
      return new Response("Unauthorized", { status: 401 });
    }
    
    if (!verifySlackSignature(body, timestamp, signature, slackSecret)) {
      console.error("Signature verification failed");
      return new Response("Unauthorized", { status: 401 });
    }

    // Parse Slack event based on content type
    let payload;
    if (contentType?.includes("application/x-www-form-urlencoded")) {
      // URL-encoded format (most Slack events)
      const params = new URLSearchParams(body);
      const payloadString = params.get("payload");
      payload = payloadString ? JSON.parse(payloadString) : Object.fromEntries(params);
    } else {
      // JSON format (URL verification)
      payload = JSON.parse(body);
    }

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

      // Process the message in the background
      EdgeRuntime.waitUntil(processSlackMessage(event, text));

      // Return immediately to acknowledge the webhook
      return new Response("OK");
    }

    return new Response("OK");
  } catch (error) {
    console.error("Error:", error);
    return new Response("Internal Server Error", { status: 500 });
  }
});
