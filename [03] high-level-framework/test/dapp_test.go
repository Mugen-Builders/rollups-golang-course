package test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/cmd/root"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository/factory"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestVotingSystem(t *testing.T) {
	suite.Run(t, new(VotingSystemSuite))
}

type VotingSystemSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *VotingSystemSuite) SetupTest() {
	repo, err := factory.NewRepositoryFromConnectionString("sqlite://:memory:")
	if err != nil {
		slog.Error("Failed to setup in-memory SQLite database", "error", err)
		os.Exit(1)
	}
	dapp := root.NewVotingSystem(repo)
	s.tester = rollmelette.NewTester(dapp)
}

func (s *VotingSystemSuite) TestAdvanceVotingHandlers() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)
	s.Contains(string(result.Notices[0].Payload), "voting created")
}

func (s *VotingSystemSuite) TestAdvanceVoterHandlers() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	createVoterInput := []byte(`{"path":"voter/create","data":{}}`)
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)
	s.Contains(string(result.Notices[0].Payload), "voter created")
}

func (s *VotingSystemSuite) TestAdvanceVotingOptionHandlers() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Nil(result.Err, "Failed to create voting")
	s.Len(result.Notices, 1, "Expected one notice for voting creation")
	s.Contains(string(result.Notices[0].Payload), "voting created")

	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Nil(result.Err, "Failed to create voting option")
	s.Len(result.Notices, 1, "Expected one notice for voting option creation")
	s.Contains(string(result.Notices[0].Payload), "voting option created")
}

func (s *VotingSystemSuite) TestDeleteVoting() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	deleteVotingInput := []byte(`{"path":"voting/delete","data":{"id":1}}`)
	result = s.tester.Advance(candidate, deleteVotingInput)
	s.Len(result.Notices, 1)
	s.Contains(string(result.Notices[0].Payload), "voting deleted")
}

func (s *VotingSystemSuite) TestDeleteVoter() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	createVoterInput := []byte(`{"path":"voter/create","data":{}}`)
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	deleteVoterInput := []byte(`{"path":"voter/delete","data":{"id":1}}`)
	result = s.tester.Advance(admin, deleteVoterInput)
	s.Len(result.Notices, 1)
	s.Contains(string(result.Notices[0].Payload), "voter deleted")
}

func (s *VotingSystemSuite) TestDeleteVotingOption() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)

	deleteOptionInput := []byte(`{"path":"voting-option/delete","data":{"id":1}}`)
	result = s.tester.Advance(candidate, deleteOptionInput)
	s.Len(result.Notices, 1)
	s.Contains(string(result.Notices[0].Payload), "voting option deleted")
}

func (s *VotingSystemSuite) TestInspectVotingHandlers() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)
	s.Equal(fmt.Sprintf(`voting created - {"id":1,"title":"Test Voting","issuer":"%s","status":"open","start_date":%d,"end_date":%d}`, candidate.Hex(), startDate, endDate), string(result.Notices[0].Payload))

	// Test FindAll
	findAllInput := []byte(`{"path":"voting","data":{}}`)
	inspectResult := s.tester.Inspect(findAllInput)
	s.Nil(inspectResult.Err)
	expectedFindAll := fmt.Sprintf(`[{"id":1,"title":"Test Voting","status":"open","start_date":%d,"end_date":%d}]`, startDate, endDate)
	s.Equal(expectedFindAll, string(inspectResult.Reports[0].Payload))

	// Test FindByID
	findByIdInput := []byte(`{"path":"voting/id","data":{"id":1}}`)
	inspectResult = s.tester.Inspect(findByIdInput)
	s.Nil(inspectResult.Err)
	expectedFindById := fmt.Sprintf(`{"id":1,"title":"Test Voting","status":"open","start_date":%d,"end_date":%d}`, startDate, endDate)
	s.Equal(expectedFindById, string(inspectResult.Reports[0].Payload))

	// Test FindAllActive
	findActiveInput := []byte(`{"path":"voting/active","data":{}}`)
	inspectResult = s.tester.Inspect(findActiveInput)
	s.Nil(inspectResult.Err)
	expectedFindActive := fmt.Sprintf(`[{"id":1,"title":"Test Voting","status":"open","start_date":%d,"end_date":%d}]`, startDate, endDate)
	s.Equal(expectedFindActive, string(inspectResult.Reports[0].Payload))

	// Test GetResults
	getResultsInput := []byte(`{"path":"voting/results","data":{"id":1}}`)
	inspectResult = s.tester.Inspect(getResultsInput)
	s.Nil(inspectResult.Err)
	expectedGetResults := `{"id":1,"title":"Test Voting","status":"open","total_votes":0,"options":[],"winner_id":0,"winner_votes":0}`
	s.Equal(expectedGetResults, string(inspectResult.Reports[0].Payload))
}

func (s *VotingSystemSuite) TestInspectVoterHandlers() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	createVoterInput := []byte(`{"path":"voter/create","data":{}}`)
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)
	s.Equal(fmt.Sprintf(`voter created - {"id":1,"address":"%s"}`, admin.Hex()), string(result.Notices[0].Payload))

	findByIdInput := []byte(`{"path":"voter/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findByIdInput)
	s.Nil(inspectResult.Err)
	expectedFindById := fmt.Sprintf(`{"id":1,"address":"%s"}`, admin.Hex())
	s.Equal(expectedFindById, string(inspectResult.Reports[0].Payload))

	findByAddressInput := []byte(fmt.Sprintf(`{"path":"voter/address","data":{"address":"%s"}}`, admin))
	inspectResult = s.tester.Inspect(findByAddressInput)
	s.Nil(inspectResult.Err)
	expectedFindByAddress := fmt.Sprintf(`{"id":1,"address":"%s"}`, admin.Hex())
	s.Equal(expectedFindByAddress, string(inspectResult.Reports[0].Payload))
}

func (s *VotingSystemSuite) TestInspectVotingOptionHandlers() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)
	s.Equal(fmt.Sprintf(`voting created - {"id":1,"title":"Test Voting","issuer":"%s","status":"open","start_date":%d,"end_date":%d}`, candidate.Hex(), startDate, endDate), string(result.Notices[0].Payload))

	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)
	s.Equal(`voting option created - {"id":1,"voting_id":1}`, string(result.Notices[0].Payload))

	findByIdInput := []byte(`{"path":"voting-option/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findByIdInput)
	s.Nil(inspectResult.Err)
	expectedFindById := `{"id":1,"voting_id":1,"vote_count":0}`
	s.Equal(expectedFindById, string(inspectResult.Reports[0].Payload))

	findByVotingIdInput := []byte(`{"path":"voting-option/voting","data":{"voting_id":1}}`)
	inspectResult = s.tester.Inspect(findByVotingIdInput)
	s.Nil(inspectResult.Err)
	expectedFindByVotingId := `[{"id":1,"voting_id":1,"vote_count":0}]`
	s.Equal(expectedFindByVotingId, string(inspectResult.Reports[0].Payload))
}

func (s *VotingSystemSuite) TestVotingWorkflow() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)
	s.Equal(fmt.Sprintf(`voting created - {"id":1,"title":"Test Voting","issuer":"%s","status":"open","start_date":%d,"end_date":%d}`, candidate.Hex(), startDate, endDate), string(result.Notices[0].Payload))

	createVoterInput := []byte(`{"path":"voter/create","data":{}}`)
	result = s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)
	s.Equal(fmt.Sprintf(`voter created - {"id":1,"address":"%s"}`, admin.Hex()), string(result.Notices[0].Payload))

	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)
	s.Equal(`voting option created - {"id":1,"voting_id":1}`, string(result.Notices[0].Payload))

	findVotingInput := []byte(`{"path":"voting/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findVotingInput)
	s.Nil(inspectResult.Err)
	expectedFindVoting := fmt.Sprintf(`{"id":1,"title":"Test Voting","status":"open","start_date":%d,"end_date":%d}`, startDate, endDate)
	s.Equal(expectedFindVoting, string(inspectResult.Reports[0].Payload))

	findVoterInput := []byte(fmt.Sprintf(`{"path":"voter/address","data":{"address":"%s"}}`, admin))
	inspectResult = s.tester.Inspect(findVoterInput)
	s.Nil(inspectResult.Err)
	expectedFindVoter := fmt.Sprintf(`{"id":1,"address":"%s"}`, admin.Hex())
	s.Equal(expectedFindVoter, string(inspectResult.Reports[0].Payload))

	findOptionInput := []byte(`{"path":"voting-option/id","data":{"id":1}}`)
	inspectResult = s.tester.Inspect(findOptionInput)
	s.Nil(inspectResult.Err)
	expectedFindOption := `{"id":1,"voting_id":1,"vote_count":0}`
	s.Equal(expectedFindOption, string(inspectResult.Reports[0].Payload))
}

func (s *VotingSystemSuite) TestInvalidPayloads() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	invalidVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, invalidVotingInput)
	s.NotNil(result.Err)

	invalidOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":0}}`)
	result = s.tester.Advance(candidate, invalidOptionInput)
	s.NotNil(result.Err)
}

func (s *VotingSystemSuite) TestDuplicateEntries() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	createVoterInput := []byte(`{"path":"voter/create","data":{}}`)
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	result = s.tester.Advance(admin, createVoterInput)
	s.NotNil(result.Err)
}

func (s *VotingSystemSuite) TestNonExistentEntities() {
	findVotingInput := []byte(`{"path":"voting/id","data":{"id":999}}`)
	inspectResult := s.tester.Inspect(findVotingInput)
	s.NotNil(inspectResult.Err)

	findVoterInput := []byte(`{"path":"voter/address","data":{"address":"0x0000000000000000000000000000000000009999"}}`)
	inspectResult = s.tester.Inspect(findVoterInput)
	s.NotNil(inspectResult.Err)

	findOptionInput := []byte(`{"path":"voting-option/id","data":{"id":999}}`)
	inspectResult = s.tester.Inspect(findOptionInput)
	s.NotNil(inspectResult.Err)
}

func (s *VotingSystemSuite) TestVotingFlow() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	baseTime := time.Now().Unix()
	startDate := baseTime + 60
	endDate := baseTime + 120

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","start_date":%d,"end_date":%d}}`, startDate, endDate))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Nil(result.Err, "Failed to create voting")
	s.Len(result.Notices, 1, "Expected one notice for voting creation")
	s.Contains(string(result.Notices[0].Payload), "voting created")

	createVoterInput := []byte(`{"path":"voter/create","data":{}}`)
	result = s.tester.Advance(admin, createVoterInput)
	s.Nil(result.Err, "Failed to create voter")
	s.Len(result.Notices, 1, "Expected one notice for voter creation")
	s.Contains(string(result.Notices[0].Payload), "voter created")

	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Nil(result.Err, "Failed to create voting option")
	s.Len(result.Notices, 1, "Expected one notice for voting option creation")
	s.Contains(string(result.Notices[0].Payload), "voting option created")

	voteInput := []byte(`{"path":"voting/vote","data":{"voting_id":1,"option_id":1}}`)
	result = s.tester.Advance(admin, voteInput)
	s.Nil(result.Err, "Failed to vote")
	s.Len(result.Notices, 1, "Expected one notice for vote")
	s.Contains(string(result.Notices[0].Payload), "vote registered")

	result = s.tester.Advance(admin, voteInput)
	s.NotNil(result.Err, "Expected error when voting twice")
}
