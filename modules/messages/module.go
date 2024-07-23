package messages

import (
	"github.com/zenrocklabs/juno/database"
	"github.com/zenrocklabs/juno/modules"
	"github.com/zenrocklabs/juno/types"
)

var _ modules.Module = &Module{}

// Module represents the module allowing to store messages properly inside a dedicated table
type Module struct {
	parser MessageAddressesParser

	db database.Database
}

func NewModule(parser MessageAddressesParser, db database.Database) *Module {
	return &Module{
		parser: parser,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "messages"
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg types.Message, tx *types.Transaction) error {
	return HandleMsg(index, msg, tx, m.parser, m.db)
}
