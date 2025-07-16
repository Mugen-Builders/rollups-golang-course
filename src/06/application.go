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

var (
	emergencyWithdrawAddress = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa") // TODO: replace with the actual address
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
		Path string          `json:"path" validate:"required"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return err
	}

	if deposit != nil {
		switch d := deposit.(type) {
		case *rollmelette.ERC20Deposit:
			if input.Path == "safe_transfer" {
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

				halfAmount := new(big.Int).Div(d.Value, big.NewInt(2))

				delegateCallVoucher, err := abiInterface.Pack("safeTransfer", d.Token, d.Sender, halfAmount)
				if err != nil {
					return err
				}

				delegateCallVoucherTargeted, err := abiInterface.Pack("safeTransferTargeted", d.Token, anyone, d.Sender, halfAmount)
				if err != nil {
					return err
				}

				env.SetERC20Balance(d.Token, d.Sender, new(big.Int).Sub(env.ERC20BalanceOf(d.Token, d.Sender), d.Value))

				env.DelegateCallVoucher(safeERC20TransferAddress, delegateCallVoucher)
				env.DelegateCallVoucher(safeERC20TransferAddress, delegateCallVoucherTargeted)
				return nil
			}
		default:
			env.Report([]byte(fmt.Sprintf("Unknown deposit type: %T", d)))
			return nil
		}
	}

	switch input.Path {
	case "emergency_erc20_withdraw":
		var emergencyInput struct {
			Token common.Address `json:"token"`
			To    common.Address `json:"to"`
		}
		if err := json.Unmarshal(input.Data, &emergencyInput); err != nil {
			return err
		}
		abiJSON := `[{
			"type":"function",
			"name":"emergencyERC20Withdraw",
			"inputs":[
				{"type":"address"},
				{"type":"address"}
			]
		}]`
		abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return err
		}
		delegateCallVoucher, err := abiInterface.Pack("emergencyERC20Withdraw", emergencyInput.Token, emergencyInput.To)
		if err != nil {
			return err
		}
		env.DelegateCallVoucher(emergencyWithdrawAddress, delegateCallVoucher)
		return nil

	case "emergency_eth_withdraw":
		var emergencyInput struct {
			To common.Address `json:"to"`
		}
		if err := json.Unmarshal(input.Data, &emergencyInput); err != nil {
			return err
		}
		abiJSON := `[{
			"type":"function",
			"name":"emergencyETHWithdraw",
			"inputs":[
				{"type":"address"}
			]
		}]`
		abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return err
		}
		delegateCallVoucher, err := abiInterface.Pack("emergencyETHWithdraw", emergencyInput.To)
		if err != nil {
			return err
		}
		env.DelegateCallVoucher(emergencyWithdrawAddress, delegateCallVoucher)
		return nil
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
