package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
)

type Application struct{}

func (a *Application) mintNFT(to common.Address, uri string) ([]byte, error) {
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

func (a *Application) deployContract(bytecode, encodedArgs []byte) ([]byte, error) {
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

	encoded := append(bytecode, encodedArgs...)

	voucher, err := abiInterface.Pack("deploy", encoded)
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

	switch input.Path {
	case "mintNFT":
		var mintNFTInput struct {
			Token common.Address `json:"token"`
			To    common.Address `json:"to"`
			URI   string         `json:"uri"`
		}
		if err := json.Unmarshal(input.Data, &mintNFTInput); err != nil {
			return err
		}
		voucher, err := a.mintNFT(mintNFTInput.To, mintNFTInput.URI)
		if err != nil {
			return err
		}
		env.Voucher(mintNFTInput.Token, big.NewInt(0), voucher)
	case "deployContract":
		var deployContractInput struct {
			ProxyDeployer common.Address `json:"proxy_deployer"`
			Bytecode      []byte         `json:"bytecode"`
			EncodedArgs   []byte         `json:"encoded_args"`
		}
		if err := json.Unmarshal(input.Data, &deployContractInput); err != nil {
			return err
		}
		voucher, err := a.deployContract(deployContractInput.Bytecode, deployContractInput.EncodedArgs)
		if err != nil {
			return err
		}
		env.Voucher(deployContractInput.ProxyDeployer, big.NewInt(0), voucher)
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
