## <div align="center">References & Instructions</div>

**Step 1:** Compile/build the application:
```bash
cartesi build
```

**Step 2:** Start the local infrastructure:
```bash
cartesi rollups start --services graphql
```

**Step 3:** Create a new To-Do item:
```bash
cartesi send generic --input='{"path":"createToDo","payload":{"title":"Create an application","description":"Use the Cartesi CLI"}}' --input-encoding=string
```

> [!NOTE]
> Replace `<application>` with your application address (e.g., `0x9321e0dd59bad3ff98836bb83403e1598a0a4478`)

**Step 4:** View and decode all outputs:
```bash
cast rpc --raw --rpc-url http://127.0.0.1:8080/rpc cartesi_listOutputs \
  '{"application":"<application>","limit":1,"offset":0,"decode":true}' \
| jq -r '.data[]?.decoded_data.payload' \
| while read -r hex; do
    if [ "$hex" != "null" ]; then
        echo "$hex" | sed 's/^0x//' | xxd -r -p
        echo
    fi
done
```

**Alternative with GraphQL:**
```bash
curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { notices(last: 1) { edges { node { payload application { name } } } } }"}' \
| jq -r '
  .data.notices.edges[0].node.payload,
  .data.notices.edges[0].node.application.name
' | {
  read hex
  read appname
  echo "Application Name: $appname"
  echo -n "Notice Payload: "
  echo "$hex" | sed 's/^0x//' | xxd -r -p
  echo
}
```

**Step 5:** Inspect all To-Dos (raw output via `jq`):
```bash
curl -X POST http://localhost:8080/inspect/<application> \
    -H "Content-Type: application/json" | jq
```


**Step 6:** Inspect all To-Dos (decoded payloads):
```bash
curl -X POST http://localhost:8080/inspect/<application> \
    -H "Content-Type: application/json" \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

**Step 7:** Update an existing To-Do item:
```bash
cartesi send generic --input='{"path":"updateToDo","payload":{"id":1,"title":"Create an application","description":"Use the Cartesi CLI","completed":true}}' --input-encoding=string
```

> [!NOTE]
> Replace `<application>` with your application address (e.g., `0x9321e0dd59bad3ff98836bb83403e1598a0a4478`)

**Step 8:** View and decode all outputs again to confirm the update:
```bash
cast rpc --raw --rpc-url http://127.0.0.1:8080/rpc cartesi_listOutputs \
  '{"application":"<application>","limit":2,"offset":0,"decode":true}' \
| jq -r '.data[]?.decoded_data.payload' \
| while read -r hex; do
    if [ "$hex" != "null" ]; then
        echo "$hex" | sed 's/^0x//' | xxd -r -p
        echo
    fi
done
```

**Alternative with GraphQL:**
```bash
curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { notices(last: 1) { edges { node { payload application { name } } } } }"}' \
| jq -r '
  .data.notices.edges[0].node.payload,
  .data.notices.edges[0].node.application.name
' | {
  read hex
  read appname
  echo "Application Name: $appname"
  echo -n "Notice Payload: "
  echo "$hex" | sed 's/^0x//' | xxd -r -p
  echo
}
```

**Step 9:** Inspect all To-Dos (decoded payloads):
```bash
curl -X POST http://localhost:8080/inspect/<application> \
    -H "Content-Type: application/json" \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

**Step 10:** Delete a To-Do item:
```bash
cartesi send generic --input='{"path":"deleteToDo","payload":{"id":1}}' --input-encoding=string
```

> [!NOTE]
> Replace `<application>` with your application address (e.g., `0x9321e0dd59bad3ff98836bb83403e1598a0a4478`)

**Step 11:** Final check — view and decode all current outputs:
```bash
cast rpc --raw --rpc-url http://127.0.0.1:8080/rpc cartesi_listOutputs \
  '{"application":"<application>","limit":2,"offset":0,"decode":true}' \
| jq -r '.data[]?.decoded_data.payload' \
| while read -r hex; do
    if [ "$hex" != "null" ]; then
        echo "$hex" | sed 's/^0x//' | xxd -r -p
        echo
    fi
done
```

**Alternative with GraphQL:**
```bash
curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { notices(last: 1) { edges { node { payload application { name } } } } }"}' \
| jq -r '
  .data.notices.edges[0].node.payload,
  .data.notices.edges[0].node.application.name
' | {
  read hex
  read appname
  echo "Application Name: $appname"
  echo -n "Notice Payload: "
  echo "$hex" | sed 's/^0x//' | xxd -r -p
  echo
}
```

**Step 12:** Stop the local infrastructure:
```bash
cartesi rollups stop
```
