package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/cmd/root"
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
	dapp := root.NewVotingSystem()
	s.tester = rollmelette.NewTester(dapp)
}

// Voting Group Tests
func (s *VotingSystemSuite) TestCreateVoting() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// Create voting
	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)
}

func (s *VotingSystemSuite) TestDeleteVoting() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	deleteVotingInput := []byte(`{"path":"voting/delete","data":{"id":1}}`)
	result := s.tester.Advance(candidate, deleteVotingInput)
	s.Len(result.Notices, 1)
}

func (s *VotingSystemSuite) TestFindAllVotings() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	// First create a voting
	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Now inspect all votings
	findAllInput := []byte(`{"path":"voting","data":{}}`)
	inspectResult := s.tester.Inspect(findAllInput)
	s.Nil(inspectResult.Err)
}

func (s *VotingSystemSuite) TestFindVotingByID() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	// First create a voting
	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Now inspect the voting by ID
	findByIdInput := []byte(`{"path":"voting/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findByIdInput)
	s.Nil(inspectResult.Err)
}

func (s *VotingSystemSuite) TestFindAllActiveVotings() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	// First create an active voting
	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Now inspect active votings
	findActiveInput := []byte(`{"path":"voting/active","data":{}}`)
	inspectResult := s.tester.Inspect(findActiveInput)
	s.Nil(inspectResult.Err)
}

func (s *VotingSystemSuite) TestGetVotingResults() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	// First create a voting
	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Create a voter
	createVoterInput := []byte(fmt.Sprintf(`{"path":"voter/create","data":{"address":"%s","name":"Test Voter"}}`, voter))
	result = s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	// Create a voting option
	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1,"title":"Test Option","description":"Test Option Description"}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)

	// Now get voting results
	getResultsInput := []byte(`{"path":"voting/results","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(getResultsInput)
	s.Nil(inspectResult.Err)
}

// Voter Group Tests
func (s *VotingSystemSuite) TestCreateVoter() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	createVoterInput := []byte(fmt.Sprintf(`{"path":"voter/create","data":{"address":"%s","name":"Test Voter"}}`, voter))
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)
}

func (s *VotingSystemSuite) TestDeleteVoter() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	deleteVoterInput := []byte(fmt.Sprintf(`{"path":"voter/delete","data":{"address":"%s"}}`, voter))
	result := s.tester.Advance(admin, deleteVoterInput)
	s.Len(result.Notices, 1)
}

func (s *VotingSystemSuite) TestFindVoterByID() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	// First create a voter
	createVoterInput := []byte(fmt.Sprintf(`{"path":"voter/create","data":{"address":"%s","name":"Test Voter"}}`, voter))
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	// Now find voter by ID
	findByIdInput := []byte(`{"path":"voter/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findByIdInput)
	s.Nil(inspectResult.Err)
}

func (s *VotingSystemSuite) TestFindVoterByAddress() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	// First create a voter
	createVoterInput := []byte(fmt.Sprintf(`{"path":"voter/create","data":{"address":"%s","name":"Test Voter"}}`, voter))
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	// Now find voter by address
	findByAddressInput := []byte(fmt.Sprintf(`{"path":"voter/address","data":{"address":"%s"}}`, voter))
	inspectResult := s.tester.Inspect(findByAddressInput)
	s.Nil(inspectResult.Err)
}

// VotingOption Group Tests
func (s *VotingSystemSuite) TestCreateVotingOption() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1,"title":"Test Option","description":"Test Option Description"}}`)
	result := s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)
}

func (s *VotingSystemSuite) TestDeleteVotingOption() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	deleteOptionInput := []byte(`{"path":"voting-option/delete","data":{"id":1}}`)
	result := s.tester.Advance(candidate, deleteOptionInput)
	s.Len(result.Notices, 1)
}

func (s *VotingSystemSuite) TestFindVotingOptionByID() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	// First create a voting
	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Create a voting option
	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1,"title":"Test Option","description":"Test Option Description"}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)

	// Now find option by ID
	findByIdInput := []byte(`{"path":"voting-option/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findByIdInput)
	s.Nil(inspectResult.Err)
}

func (s *VotingSystemSuite) TestFindAllOptionsByVotingID() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	// First create a voting
	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Create a voting option
	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1,"title":"Test Option","description":"Test Option Description"}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)

	// Now find all options for the voting
	findByVotingIdInput := []byte(`{"path":"voting-option/voting","data":{"voting_id":1}}`)
	inspectResult := s.tester.Inspect(findByVotingIdInput)
	s.Nil(inspectResult.Err)
}

// Integration Tests
func (s *VotingSystemSuite) TestVotingWorkflow() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// Create voting
	createVotingInput := []byte(fmt.Sprintf(`{"path":"voting/create","data":{"title":"Test Voting","description":"Test Description","closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result := s.tester.Advance(candidate, createVotingInput)
	s.Len(result.Notices, 1)

	// Create voter
	createVoterInput := []byte(fmt.Sprintf(`{"path":"voter/create","data":{"address":"%s","name":"Test Voter"}}`, voter))
	result = s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	// Create voting option
	createOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":1,"title":"Test Option","description":"Test Option Description"}}`)
	result = s.tester.Advance(candidate, createOptionInput)
	s.Len(result.Notices, 1)

	// Verify voting was created
	findVotingInput := []byte(`{"path":"voting/id","data":{"id":1}}`)
	inspectResult := s.tester.Inspect(findVotingInput)
	s.Nil(inspectResult.Err)

	// Verify voter was created
	findVoterInput := []byte(fmt.Sprintf(`{"path":"voter/address","data":{"address":"%s"}}`, voter))
	inspectResult = s.tester.Inspect(findVoterInput)
	s.Nil(inspectResult.Err)

	// Verify option was created
	findOptionInput := []byte(`{"path":"voting-option/id","data":{"id":1}}`)
	inspectResult = s.tester.Inspect(findOptionInput)
	s.Nil(inspectResult.Err)
}

// Validation Tests
func (s *VotingSystemSuite) TestInvalidPayloads() {
	candidate := common.HexToAddress("0x0000000000000000000000000000000000000007")

	// Test invalid voting creation
	invalidVotingInput := []byte(`{"path":"voting/create","data":{"title":"","description":""}}`)
	result := s.tester.Advance(candidate, invalidVotingInput)
	s.NotNil(result.Err)

	// Test invalid voter creation
	invalidVoterInput := []byte(`{"path":"voter/create","data":{"address":"","name":""}}`)
	result = s.tester.Advance(candidate, invalidVoterInput)
	s.NotNil(result.Err)

	// Test invalid option creation
	invalidOptionInput := []byte(`{"path":"voting-option/create","data":{"voting_id":0,"title":"","description":""}}`)
	result = s.tester.Advance(candidate, invalidOptionInput)
	s.NotNil(result.Err)
}

func (s *VotingSystemSuite) TestDuplicateEntries() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	voter := common.HexToAddress("0x0000000000000000000000000000000000000001")

	// Create voter first time
	createVoterInput := []byte(fmt.Sprintf(`{"path":"voter/create","data":{"address":"%s","name":"Test Voter"}}`, voter))
	result := s.tester.Advance(admin, createVoterInput)
	s.Len(result.Notices, 1)

	// Try to create same voter again
	result = s.tester.Advance(admin, createVoterInput)
	s.NotNil(result.Err)
}

func (s *VotingSystemSuite) TestNonExistentEntities() {
	// Try to find non-existent voting
	findVotingInput := []byte(`{"path":"voting/id","data":{"id":999}}`)
	inspectResult := s.tester.Inspect(findVotingInput)
	s.NotNil(inspectResult.Err)

	// Try to find non-existent voter
	findVoterInput := []byte(`{"path":"voter/address","data":{"address":"0x0000000000000000000000000000000000009999"}}`)
	inspectResult = s.tester.Inspect(findVoterInput)
	s.NotNil(inspectResult.Err)

	// Try to find non-existent option
	findOptionInput := []byte(`{"path":"voting-option/id","data":{"id":999}}`)
	inspectResult = s.tester.Inspect(findOptionInput)
	s.NotNil(inspectResult.Err)
}
