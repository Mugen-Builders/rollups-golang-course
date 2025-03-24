#!/bin/bash

stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

RPC_URL="http://localhost:8545"
DAPP_ADDRESS="0x88cf12bbe0d2e20748fa6a690995fc900e45c195"

sendInput() {
    local payload="$1"
    local hexPayload
    hexPayload=$(stringToHex "$payload")

    echo "Sending input: $hexPayload"
    cartesi send generic \
        --input="$hexPayload" \
        --input-encoding=hex \
        --dapp="$DAPP_ADDRESS" \
        --chain-id=31337 \
        --mnemonic-index=0 \
        --rpc-url="$RPC_URL" \
        --mnemonic="test test test test test test test test test test test junk"
}

# Edit the payload before sending it.
echo "Creating To-Do..."
sendInput '{"path":"deleteToDo","payload":{"id":1}}'
sleep 1