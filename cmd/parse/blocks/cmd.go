package blocks

import (
	"github.com/spf13/cobra"

	parsecmdtypes "github.com/zenrocklabs/juno/cmd/parse/types"
)

// NewBlocksCmd returns the Cobra command that allows to fix all the things related to blocks
func NewBlocksCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocks",
		Short: "Fix things related to blocks and transactions",
	}

	cmd.AddCommand(
		newAllCmd(parseConfig),
		newMissingCmd(parseConfig),
	)

	return cmd
}
