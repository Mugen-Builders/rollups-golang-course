package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type Application struct{}

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

func (a *Application) DeployContract(bytecode []byte) ([]byte, error) {
	abiJSON := `[{
		"type":"function",
		"name":"deploy",
		"inputs":[
			{"type":"bytes"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	voucher, err := abiInterface.Pack("deploy", bytecode)
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
	case "mint_nft":
		var mintNFTInput struct {
			Token common.Address `json:"token"`
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
		env.Voucher(mintNFTInput.Token, big.NewInt(0), voucher)
	case "deploy_contract":
		var deployContractInput struct {
			Deployer common.Address `json:"deployer"`
			Bytecode []byte         `json:"bytecode"`
		}
		if err := json.Unmarshal(input.Data, &deployContractInput); err != nil {
			return err
		}
		voucher, err := a.DeployContract(deployContractInput.Bytecode)
		if err != nil {
			return err
		}
		env.Voucher(deployContractInput.Deployer, big.NewInt(0), voucher)
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
