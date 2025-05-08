#!/bin/bash

# API 키 설정
API_KEY="YOUR_API_KEY"
API_SECRET="YOUR_SECRECT_KEY"

# 기본 설정
TIMESTAMP=$(($(gdate +%s%3N))) # 밀리초
RECV_WINDOW=5000

# 계정 파라미터
CATEGORY="spot"
SYMBOL="BTCUSDT"
SIDE="Buy"          # Buy 또는 Sell
ORDER_TYPE="Market" # Market 또는 Limit
QTY="10"
# QTY="0.001"
PRICE="50000" # Limit 주문일 경우에만 필요

# API 엔드포인트
ENDPOINT="/v5/order/create"
BASE_URL="https://api.bybit.com"
TESTNET_URL="https://api-testnet.bybit.com"
FULL_URL="${BASE_URL}${ENDPOINT}"

# Request body 생성
REQUEST_BODY="{\"category\":\"$CATEGORY\",\"symbol\":\"$SYMBOL\",\"side\":\"$SIDE\",\"orderType\":\"$ORDER_TYPE\",\"qty\":\"$QTY\"}"
# REQUEST_BODY="{\"category\":\"$CATEGORY\",\"symbol\":\"$SYMBOL\",\"side\":\"$SIDE\",\"orderType\":\"$ORDER_TYPE\",\"qty\":\"$QTY\",\"price\":\"$PRICE\"}"

# Sign string 생성 (POST 요청용)
SIGN_STRING="${TIMESTAMP}${API_KEY}${RECV_WINDOW}${REQUEST_BODY}"
SIGN=$(echo -n "$SIGN_STRING" | openssl dgst -sha256 -hmac "$API_SECRET" | sed 's/^.* //')

# API 요청
RESPONSE=$(
  curl -X POST "$FULL_URL" \
    -H "X-BAPI-API-KEY: $API_KEY" \
    -H "X-BAPI-SIGN: $SIGN" \
    -H "X-BAPI-TIMESTAMP: $TIMESTAMP" \
    -H "X-BAPI-RECV-WINDOW: $RECV_WINDOW" \
    -H "Content-Type: application/json" \
    -d "$REQUEST_BODY"
)

# 결과 출력
echo "$RESPONSE" | jq
