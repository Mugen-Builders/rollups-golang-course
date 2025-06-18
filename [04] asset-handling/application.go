package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
)

type Application struct{}

func (a *Application) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	switch d := deposit.(type) {
	case *rollmelette.EtherDeposit:

		env.Notice([]byte(fmt.Sprintf("1 - Ether balance of %s: %d, before transfer to 0x0000000000000000000000000000000000000000", d.Sender, env.EtherBalanceOf(d.Sender))))
		env.EtherTransfer(d.Sender, common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"2 - Balance of %s: %d before transfer to %s",
			common.HexToAddress("0x0000000000000000000000000000000000000000").Hex(),
			env.EtherBalanceOf(common.HexToAddress("0x0000000000000000000000000000000000000000")),
			d.Sender,
		)))
		env.EtherTransfer(common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf("3 - Ether balance of %s: %d before withdraw", d.Sender, env.EtherBalanceOf(d.Sender))))
		env.EtherWithdraw(d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf("4 - Ether balance of %s: %d after withdraw", d.Sender, env.EtherBalanceOf(d.Sender))))
	case *rollmelette.ERC20Deposit:

		env.Notice([]byte(fmt.Sprintf("1 - ERC20 balance of %s: %d before transfer to 0x0000000000000000000000000000000000000000", d.Sender, env.ERC20BalanceOf(d.Token, d.Sender))))
		env.ERC20Transfer(d.Token, d.Sender, common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"2 - Balance of %s: %d before transfer to %s",
			common.HexToAddress("0x0000000000000000000000000000000000000000").Hex(),
			env.ERC20BalanceOf(d.Token, common.HexToAddress("0x0000000000000000000000000000000000000000")),
			d.Sender,
		)))
		env.ERC20Transfer(d.Token, common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf("3 - ERC20 balance of %s: %d before withdraw", d.Sender, env.ERC20BalanceOf(d.Token, d.Sender))))
		env.ERC20Withdraw(d.Token, d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf("4 - ERC20 balance of %s: %d after withdraw", d.Sender, env.ERC20BalanceOf(d.Token, d.Sender))))
	default:
		env.Report([]byte(fmt.Sprintf("Unknown deposit type: %T", d)))
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
