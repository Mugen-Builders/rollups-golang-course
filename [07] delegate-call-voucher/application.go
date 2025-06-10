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
	"github.com/rollmelette/rollmelette"
)

var (
	safeERC20TransferAddress = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa") // TODO: replace with the actual address
	anyone                   = common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")
)

type Application struct{}

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

	switch d := deposit.(type) {
	case *rollmelette.ERC20Deposit:
		if input.Path == "safe_transfer" {
			var safeTransferInput struct {
				Token  common.Address `json:"token"`
				To     common.Address `json:"to"`
				Amount *big.Int       `json:"amount"`
			}
			if err := json.Unmarshal(input.Data, &safeTransferInput); err != nil {
				return err
			}

			abiJSON := `[{
				"type":"function",
				"name":"safeTransfer",
				"inputs":[
					{"type":"address"},
					{"type":"address"},
					{"type":"uint256"}
				]
			},
			{
				"type":"function",
				"name":"safeTransferTargeted",
				"inputs":[
					{"type":"address"},
					{"type":"address"},
					{"type":"address"},
					{"type":"uint256"}
				]
			}]`
			abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
			if err != nil {
				return err
			}

			halfAmount := new(big.Int).Div(safeTransferInput.Amount, big.NewInt(2))

			delegateCallVoucher, err := abiInterface.Pack("safeTransfer", safeTransferInput.Token, safeTransferInput.To, halfAmount)
			if err != nil {
				return err
			}

			delegateCallVoucherTargeted, err := abiInterface.Pack("safeTransferTargeted", safeTransferInput.Token, anyone, safeTransferInput.To, halfAmount)
			if err != nil {
				return err
			}

			env.DelegateCallVoucher(safeERC20TransferAddress, delegateCallVoucher)
			env.DelegateCallVoucher(safeERC20TransferAddress, delegateCallVoucherTargeted)
			return nil
		} else {
			env.Report([]byte(fmt.Sprintf("Unknown path: %s", input.Path)))
			return nil
		}
	default:
		env.Report([]byte(fmt.Sprintf("Unknown deposit type: %T", d)))
		return nil
	}
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
