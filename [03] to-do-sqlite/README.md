## <div align="center">Instructions</div>

**Step 1:** Compile/build the application:
```bash
cartesi build
```

**Step 2:** Run the application:
```bash
cartesi run
```

**Step 3:** Grant execution permission to all scripts in the tools folder:
```bash
chmod +x ./tools/*.sh
```

**Step 4:** Create a new To-Do item:
```bash
./tools/create_to_do.sh
```

**Step 5:** Inspect all To-Dos (raw output via `jq`):
```bash
curl 'http://localhost:8080/inspect/todos' | jq
```

**Step 6:** Inspect all To-Dos (decoded payload):
```bash
curl 'http://localhost:8080/inspect/todos' \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

**Step 7:** Update an existing To-Do item:
```bash
./tools/update_to_do.sh
```

**Step 8:** Inspect all To-Dos (decoded payload) again to confirm changes:
```bash
curl 'http://localhost:8080/inspect/todos' \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

**Step 9:** Delete a To-Do item:
```bash
./tools/delete_to_do.sh
```

**Step 10:** Inspect all To-Dos (decoded payload) one more time:
```bash
curl 'http://localhost:8080/inspect/todos' \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

> [!ALERT]
> Proceed with the command below only when you are running the Cartesi Rollups infrastructure locally.

**Step 11:** Access the Cartesi explorer to see all details and outputs for each input submitted:
<br>


[![Docs]][Link-docs]

[Docs]: https://img.shields.io/badge/Cartesi-Explorer-79F7FA?style=for-the-badge
[Link-docs]: http://localhost:8080/explorer