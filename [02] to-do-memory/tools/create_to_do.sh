#!/bin/bash

stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

RPC_URL="http://localhost:8545"
DAPP_ADDRESS="0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"

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
        --mnemonic-passphrase="test test test test test test test test test test test junk"
}

# Edit the payload before sending it.
echo "Creating To-Do..."
sendInput '{"path":"createToDo","payload":{"title":"create a application","description":"Use the Cartesi Cli"}}'
sleep 1