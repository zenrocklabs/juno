package main

import (
	"os"

	"github.com/zenrocklabs/juno/cmd/parse/types"

	"github.com/zenrocklabs/juno/modules/messages"
	"github.com/zenrocklabs/juno/modules/registrar"

	"github.com/zenrocklabs/juno/cmd"
)

func main() {
	// JunoConfig the runner
	config := cmd.NewConfig("juno").
		WithParseConfig(types.NewConfig().
			WithRegistrar(registrar.NewDefaultRegistrar(
				messages.CosmosMessageAddressesParser,
			)),
		)

	// Run the commands and panic on any error
	exec := cmd.BuildDefaultExecutor(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
