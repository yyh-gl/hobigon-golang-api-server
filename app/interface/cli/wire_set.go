package cli

import "github.com/google/wire"

// WireSet : interface層のWireSet（CLI用）
var WireSet = wire.NewSet(
	NewNotification,
)
