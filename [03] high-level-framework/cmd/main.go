package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/cmd/root"
	"github.com/rollmelette/rollmelette"
)

func main() {
	r := root.NewVotingSystem()
	ctx := context.Background()
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(ctx, opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}
