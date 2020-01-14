package app

import "github.com/google/wire"

var APISet = wire.NewSet(
	NewAPILogger,
)

var CLISet = wire.NewSet(
	NewCLILogger,
)
