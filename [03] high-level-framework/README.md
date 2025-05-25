## <div align="center">References & Instructions</div>

**Step 1:** Compile/build the application:
```bash
cartesi build
```

**Step 2:** Start the local infrastructure:
```bash
 cartesi rollups start
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

**Step 12:** Stop the local infrastructure:
```bash
cartesi rollups stop
```


















# Voting System Use Cases

## 1. Voter Management

### Register Voter
```go
RegisterVoter(address string) error
```
- Register a new voter with their address
- Validate if address is valid
- Check if voter already exists

### Find Voter
```go
FindVoterByID(id int) (*Voter, error)
FindVoterByAddress(address string) (*Voter, error)
```
- Find a voter by ID or address
- Return error if not found

## 2. Voting Management

### Create Voting
```go
CreateVoting(title, description string, startTime, endTime time.Time) (*Voting, error)
```
- Create a new voting session
- Validate if dates are valid
- Check if voting with same title exists

### Find Voting
```go
FindVotingByID(id int) (*Voting, error)
FindAllVotings() ([]*Voting, error)
FindAllActiveVotings() ([]*Voting, error)
```
- Find a voting session by ID
- List all voting sessions
- List only active voting sessions
- Include status and results

## 3. Voting Option Management

### Create Option
```go
CreateVotingOption(votingID int, description string) (*VotingOption, error)
```
- Add a new voting option
- Validate if voting exists
- Check if option already exists

### Find Options
```go
FindAllOptionsByVotingID(votingID int) ([]*VotingOption, error)
```
- List all options for a voting session
- Include vote counts

## 4. Vote Casting

### Cast Vote
```go
CastVote(voterID, votingID, optionID int) error
```
- Register a vote
- Validate if voter exists
- Validate if voting is active
- Validate if option exists
- Check if voter already voted
- Increment vote counter

### Check Vote Status
```go
HasVoted(voterID, votingID int) (bool, error)
```
- Check if a voter has already voted
- Validate if voter exists
- Validate if voting exists

## 5. Voting Period Management

### Control Voting Period
```go
StartVoting(id int) error
EndVoting(id int) error
```
- Start a voting session
- End a voting session
- Validate if voting exists
- Check if it's active/inactive
- Calculate final results when ending

## 6. Result Management

### Count Votes
```go
CountVotes(votingID int) (map[int]int, error)
```
- Count votes by option
- Validate if voting exists
- Return map of option -> vote count

### Get Results
```go
GetVotingResults(votingID int) (*VotingResults, error)
```
- Get complete voting results
- Include total votes
- Include percentages
- Include winning option