#!/bin/bash

# ref link.
# https://platform.openai.com/docs/guides/text?api-mode=responses&lang=curl
# https://platform.openai.com/docs/api-reference/responses/create

OPENAI_API_KEY="YOUR_API_KEY"
ORG_ID="YOUR_ORG_ID"
PROJECT_ID="YOUR_PROJECT_ID"

curl -G \
    "https://api.openai.com/v1/models" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -H "OpenAI-Organization: $ORG_ID" \
    -H "OpenAI-Project: $PROJECT_ID" |
    jq '.'

GPT_MODEL="gpt-4.1"
INSTRUCTIONS="Talk like a pirate."
INPUT="Are semicolons optional in JavaScript?"

DATA=$(jq -n \
    --arg model "$GPT_MODEL" \
    --arg instructions "$INSTRUCTIONS" \
    --arg input "$INPUT" \
    '{
    model: $model,
    instructions: $instructions,
    input: $input
  }')

curl -X POST \
    "https://api.openai.com/v1/responses" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -d "$DATA" |
    jq '.'
