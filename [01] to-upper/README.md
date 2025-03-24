## <div align="center">References & Instructions</div>

**Step 1:** Compile/build the application:
```bash
cartesi build
```

**Step 2:** Start the local infrastructure:
```bash
cartesi rollups start --services graphql
```

**Step 3:** Send an input:
```bash
❯ cartesi send generic --input='cartesi is awesome!' --input-encoding=string
```

> [!NOTE]
> Replace `<application>` with your application address (e.g., `0x9321e0dd59bad3ff98836bb83403e1598a0a4478`)

**Step 4:** Inspect last transformed input (raw output via `jq`):
```bash
curl -X POST http://localhost:8080/inspect/<application> \
    -H "Content-Type: application/json" | jq
```

**Step 5:** Inspect last transformed input (decoded):
```bash
curl -X POST http://localhost:8080/inspect/<application> \
    -H "Content-Type: application/json" \
    | jq -r '.reports[0].payload' \
```

**Step 6:** Stop the local infrastructure:
```bash
cartesi rollups stop
```