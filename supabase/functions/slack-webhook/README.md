# Slack Webhook Function

This Supabase Edge Function serves as a bridge between Slack and the Personal Agent system.

## Setup

### 1. Environment Variables

Copy `.env.local.example` to `.env.local` and fill in your credentials:

```bash
cp ../env.local.example .env.local
```

Required environment variables:
- `SLACK_SIGNING_SECRET`: Your Slack app's signing secret
- `SLACK_BOT_TOKEN`: Your Slack bot's OAuth token
- `OPENAI_API_KEY`: Your OpenAI API key
- Database connection details for your existing PostgreSQL instance

### 2. Local Development

```bash
# Start Supabase locally
supabase start

# Serve the function locally
supabase functions serve slack-webhook --env-file .env.local

# Test with curl
curl -X POST http://localhost:54321/functions/v1/slack-webhook \
  -H "Content-Type: application/json" \
  -d '{"type": "url_verification", "challenge": "test_challenge"}'
```

### 3. Deployment

```bash
# Set secrets in Supabase
supabase secrets set SLACK_SIGNING_SECRET=<your-slack-signing-secret>
supabase secrets set SLACK_BOT_TOKEN=<your-slack-bot-token>
supabase secrets set OPENAI_API_KEY=<your-openai-api-key>
supabase secrets set DB_HOST=<your-db-host>
supabase secrets set DB_PORT=<your-db-port>
supabase secrets set DB_NAME=personal_agent
supabase secrets set DB_USER=<your-db-user>
supabase secrets set DB_PASSWORD=<your-db-password>

# Deploy the function
supabase functions deploy slack-webhook
```

### 4. Slack App Configuration

1. Go to your Slack App settings at https://api.slack.com/apps
2. Navigate to "Event Subscriptions"
3. Enable events and set the Request URL to:
   ```
   https://<your-project-ref>.supabase.co/functions/v1/slack-webhook
   ```
4. Subscribe to these bot events:
   - `app_mention`
   - `message.im`
5. Navigate to "OAuth & Permissions" and add these scopes:
   - `chat:write`
   - `app_mentions:read`
   - `im:history`
6. Install the app to your workspace

## Testing

Once deployed and configured, you can test by:
1. Mentioning your bot in a channel: `@YourBotName hello`
2. Sending a direct message to your bot
3. The bot will respond using the Personal Agent's knowledge base