#!/bin/bash

# API 키 설정
API_KEY="YOUR_API_KEY"
API_SECRET="YOUR_SECRECT_KEY"

# 기본 설정
TIMESTAMP=$(($(gdate +%s%3N))) # 밀리초
RECV_WINDOW=5000

# 계정 파라미터
ACCOUNT_TYPE="UNIFIED"
COIN="BTC"

# Query string 생성
QUERY="accountType=${ACCOUNT_TYPE}&coin=${COIN}"

# API 엔드포인트
ENDPOINT="/v5/account/wallet-balance"
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

# USDT 잔고만 출력
BTC_WALLET=$(echo "$RESPONSE" | jq -r '.result.list[0].coin[] | select(.coin=="BTC") | .walletBalance')
echo "BTC 지갑 잔고: $BTC_WALLET"

echo "$RESPONSE" | jq -r '
  .result.list[0].coin[]
  | select(.coin=="BTC")
  | "BTC 잔고: \(.walletBalance)\nBTC 총 평가금액: \(.equity)\nBTC USD 가치: \(.usdValue)\nBTC 누적 실현손익: \(.cumRealisedPnl)"'
