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
		env.Notice([]byte(fmt.Sprintf("Balance of %s: %d", d.Sender, env.EtherBalanceOf(d.Sender))))
		env.EtherTransfer(d.Sender, common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Value)
		env.Notice([]byte(fmt.Sprintf("Balance of %x: %d",
			common.HexToAddress("0x0000000000000000000000000000000000000000"),
			env.EtherBalanceOf(common.HexToAddress("0x0000000000000000000000000000000000000000")),
		)))
		env.EtherTransfer(common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Sender, d.Value)
		env.Notice([]byte(fmt.Sprintf("Final balance of %s: %d", d.Sender, env.EtherBalanceOf(d.Sender))))
		env.EtherWithdraw(d.Sender, d.Value)
	case *rollmelette.ERC20Deposit:
		env.Notice([]byte(fmt.Sprintf("Balance of %s: %d", d.Sender, env.ERC20BalanceOf(d.Token, d.Sender))))
		env.ERC20Transfer(d.Token, d.Sender, common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Amount)
		env.Notice([]byte(fmt.Sprintf("Balance of %x: %d",
			common.HexToAddress("0x0000000000000000000000000000000000000000"),
			env.ERC20BalanceOf(d.Token, common.HexToAddress("0x0000000000000000000000000000000000000000")),
		)))
		env.ERC20Transfer(d.Token, common.HexToAddress("0x0000000000000000000000000000000000000000"), d.Sender, d.Amount)
		env.Notice([]byte(fmt.Sprintf("Final balance of %s: %d", d.Sender, env.ERC20BalanceOf(d.Token, d.Sender))))
		env.ERC20Withdraw(d.Token, d.Sender, d.Amount)
	default:
		addressBook := rollmelette.NewAddressBook()
		switch metadata.MsgSender {
		case addressBook.ERC721Portal:
			// TODO: Handle ERC721 deposits
			env.Notice([]byte("ERC721 deposit received"))
		case addressBook.ERC1155SinglePortal:
			// TODO: Handle ERC1155 single deposits
			env.Notice([]byte("ERC1155 single deposit received"))
		case addressBook.ERC1155BatchPortal:
			// TODO: Handle ERC1155 batch deposits
			env.Notice([]byte("ERC1155 batch deposit received"))
		}
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
