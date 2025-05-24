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

## Implementation Guidelines

Each use case should:
1. Validate its inputs
2. Handle errors appropriately
3. Maintain data consistency
4. Follow business rules
5. Return clear results

## Error Handling

Common error scenarios to handle:
- Invalid input data
- Not found resources
- Duplicate entries
- Invalid state transitions
- Unauthorized operations
- System errors 