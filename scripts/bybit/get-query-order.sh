#!/bin/bash

# API 키 설정
API_KEY="YOUR_API_KEY"
API_SECRET="YOUR_SECRET_KEY"

# 기본 설정
TIMESTAMP=$(($(gdate +%s%3N))) # 밀리초
RECV_WINDOW=5000

# 조회 파라미터
CATEGORY="spot" # spot, linear, inverse
SYMBOL="BTCUSDT"
ORDER_ID="1943731665738492160"      # 조회할 주문 ID (선택사항)
ORDER_LINK_ID="1943731665755269376" # 조회할 주문 링크 ID (선택사항)

# Query string 생성
QUERY="category=${CATEGORY}&symbol=${SYMBOL}"

# orderId가 있으면 추가
if [ ! -z "$ORDER_ID" ]; then
  QUERY="${QUERY}&orderId=${ORDER_ID}"
fi

# orderLinkId가 있으면 추가
if [ ! -z "$ORDER_LINK_ID" ]; then
  QUERY="${QUERY}&orderLinkId=${ORDER_LINK_ID}"
fi

# API 엔드포인트
ENDPOINT="/v5/order/realtime"
BASE_URL="https://api.bybit.com"
TESTNET_URL="https://api-testnet.bybit.com"
FULL_URL="${BASE_URL}${ENDPOINT}?${QUERY}"

# Sign string: GET에서는 (timestamp + apiKey + recvWindow + query string)
SIGN_STRING="${TIMESTAMP}${API_KEY}${RECV_WINDOW}${QUERY}"
SIGN=$(echo -n "$SIGN_STRING" | openssl dgst -sha256 -hmac "$API_SECRET" | sed 's/^.* //')

# API 요청
RESPONSE=$(curl -s -X GET "$FULL_URL" \
  -H "X-BAPI-API-KEY: $API_KEY" \
  -H "X-BAPI-SIGN: $SIGN" \
  -H "X-BAPI-TIMESTAMP: $TIMESTAMP" \
  -H "X-BAPI-RECV-WINDOW: $RECV_WINDOW")

# 결과 출력
echo "$RESPONSE" | jq

# 주문 상태만 출력
echo "$RESPONSE" | jq -r '.result.list[] | "주문 ID: \(.orderId)\n주문 링크 ID: \(.orderLinkId)\n상태: \(.orderStatus)\n수량: \(.qty)\n체결 수량: \(.cumExecQty)\n가격: \(.price)\n생성 시간: \(.createdTime)"'
