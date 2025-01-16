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

**Step 4:** Create a new todo item:
```bash
./tools/create_todo.sh
```

**Step 5:** Inspect all todos (raw output via `jq`):
```bash
curl 'http://localhost:8080/inspect/todos' | jq
```

**Step 6:** Inspect all todos (decoded payload):
```bash
curl 'http://localhost:8080/inspect/todos' \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

**Step 7:** Update an existing todo item:
```bash
./tools/update_todo.sh
```

**Step 8:** Inspect all todos (decoded payload) again to confirm changes:
```bash
curl 'http://localhost:8080/inspect/todos' \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```

**Step 9:** Delete a todo item:
```bash
./tools/delete_todo.sh
```

**Step 10:** Inspect all todos (decoded payload) one more time:
```bash
curl 'http://localhost:8080/inspect/todos' \
    | jq -r '.reports[0].payload' \
    | sed 's/^0x//' \
    | xxd -r -p \
    | jq
```