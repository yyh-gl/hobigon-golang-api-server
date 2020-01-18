package rest

import "github.com/google/wire"

// WireSet : interface層のWireSet（API用）
var WireSet = wire.NewSet(
	NewNotificationHandler,
	NewBirthday,
	NewBlog,
)
