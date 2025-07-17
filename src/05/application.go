package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

var (
	NFTFactoryAddress = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa") // TODO: replace with the actual address
)

type Application struct {
	NFT common.Address
}

func (a *Application) DeployNFT(initialOwner common.Address, salt common.Hash) ([]byte, error) {
	abiJSON := `[{
		"type":"function",
		"name":"newNFT",
		"inputs":[
			{"type":"address"},
			{"type":"bytes32"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	voucher, err := abiInterface.Pack("newNFT", initialOwner, salt)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (a *Application) MintNFT(to common.Address, uri string) ([]byte, error) {
	abiJSON := `[{
		"type":"function",
		"name":"safeMint",
		"inputs":[
			{"type":"address"},
			{"type":"string"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	voucher, err := abiInterface.Pack("safeMint", to, uri)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (a *Application) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	var input struct {
		Path string          `json:"path"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	switch input.Path {
	case "deploy_nft":
		bytecode, err := getNFTBytecode()
		if err != nil {
			return err
		}

		addressType, _ := abi.NewType("address", "", nil)
		constructorArgs, err := abi.Arguments{
			{Type: addressType},
		}.Pack(metadata.AppContract)
		if err != nil {
			return fmt.Errorf("error encoding constructor args: %w", err)
		}

		a.NFT = crypto.CreateAddress2(
			NFTFactoryAddress,
			common.HexToHash(strconv.Itoa(metadata.Index)),
			crypto.Keccak256(append(bytecode, constructorArgs...)),
		)

		deployNFTPayload, err := a.DeployNFT(
			metadata.AppContract,
			common.HexToHash(strconv.Itoa(metadata.Index)),
		)
		if err != nil {
			return err
		}
		env.Voucher(NFTFactoryAddress, big.NewInt(0), deployNFTPayload)
	case "mint_nft":
		var mintNFTInput struct {
			To    common.Address `json:"to"`
			URI   string         `json:"uri"`
		}
		if err := json.Unmarshal(input.Data, &mintNFTInput); err != nil {
			return err
		}
		voucher, err := a.MintNFT(mintNFTInput.To, mintNFTInput.URI)
		if err != nil {
			return err
		}
		env.Voucher(a.NFT, big.NewInt(0), voucher)
	default:
		// using report to log an error
		env.Report([]byte(fmt.Sprintf("Unknown path: %s", input.Path)))
		return fmt.Errorf("unknown path: %s", input.Path)
	}
	return nil
}

func (a *Application) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	return nil
}

func main() {
	ctx := context.Background()
	opts := rollmelette.NewRunOpts()
	app := new(Application)
	err := rollmelette.Run(ctx, opts, app)
	if err != nil {
		slog.Error("application error", "error", err)
	}
}
