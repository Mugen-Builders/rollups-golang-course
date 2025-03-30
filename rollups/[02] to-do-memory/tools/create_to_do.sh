#!/bin/bash

stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

RPC_URL="http://localhost:8545"
DAPP_ADDRESS="0xf879c61f11639116fa1d932aacaaaf97c4f56b8a"

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
        --rpc-url="http://127.0.0.1:8545" \
        --mnemonic="test test test test test test test test test test test junk"
}

# Edit the payload before sending it.
echo "Creating To-Do..."
sendInput '{"path":"createToDo","payload":{"title":"create a application","description":"Use the Cartesi Cli"}}'
sleep 1