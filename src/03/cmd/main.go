package main

import (
	"os"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/cmd/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
