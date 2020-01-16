package app

import "github.com/google/wire"

// APISet : アプリ共通系のWireSet（API用）
var APISet = wire.NewSet(
	NewAPILogger,
)

// CLISet : アプリ共通系のWireSet（CLI用）
var CLISet = wire.NewSet(
	NewCLILogger,
)
