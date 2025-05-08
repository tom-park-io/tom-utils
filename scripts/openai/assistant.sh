#!/bin/bash

# ref link.
# https://platform.openai.com/docs/assistants/quickstart

OPENAI_API_KEY="YOUR_API_KEY"
OPENAI_ASSISTANT_ID="YOUR_ASSISTANT_ID"

ORG_ID="YOUR_ORG_ID"
PROJECT_ID="YOUR_PROJECT_ID"

# 0. Fetch Assistants
curl -s \
  https://api.openai.com/v1/assistants \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "OpenAI-Beta: assistants=v2" |
  jq '.'

# 0-1. Fetch Assistant
curl -s \
  https://api.openai.com/v1/assistants/$OPENAI_ASSISTANT_ID \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "OpenAI-Beta: assistants=v2" |
  jq '.'

# 1. Create a thread
THREAD_ID=$(curl -s -X POST \
  "https://api.openai.com/v1/threads" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -H "OpenAI-Beta: assistants=v2" | jq -r '.id')
echo "🧵 Created thread: $THREAD_ID"

# 2. Add a message to the thread (role: user)
# Input values
ROLE="user" # user, assistant, or tool
CONTENT="I need to solve the equation 3x + 11 = 14. Can you help me?"

# JSON payload for message
DATA=$(jq -n \
  --arg role "$ROLE" \
  --arg content "$CONTENT" \
  '{
    role: $role,
    content: $content
  }')

# Send message to thread
curl -s -X POST \
  "https://api.openai.com/v1/threads/$THREAD_ID/messages" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "OpenAI-Beta: assistants=v2" \
  -d "$DATA" >/dev/null
echo "📩 Message added to thread"

# 3. Run the assistant
RUN_ID=$(curl -s -X POST \
  "https://api.openai.com/v1/threads/$THREAD_ID/runs" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "OpenAI-Beta: assistants=v2" \
  -d "{\"assistant_id\": \"$OPENAI_ASSISTANT_ID\"}" | jq -r '.id')
echo "🏃 Started run: $RUN_ID"

# 4. Poll until the run is complete
STATUS="in_progress"
while [[ "$STATUS" != "completed" && "$STATUS" != "failed" && "$STATUS" != "cancelled" ]]; do
  sleep 1
  STATUS=$(curl -s \
    "https://api.openai.com/v1/threads/$THREAD_ID/runs/$RUN_ID" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -H "OpenAI-Beta: assistants=v2" |
    jq -r '.status')
  echo "⏳ Waiting for run... ($STATUS)"
done

# 5. Get the final messages from the thread
echo "✅ Run completed. Retrieving messages..."
curl -s https://api.openai.com/v1/threads/$THREAD_ID/messages \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "OpenAI-Beta: assistants=v2" |
  jq -r '.data[] | select(.role=="assistant") | .content[0].text.value'
