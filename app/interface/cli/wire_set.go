package cli

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewNotificationHandler,
)
