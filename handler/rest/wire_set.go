package rest

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewNotificationHandler,
	NewBirthdayHandler,
	NewBlogHandler,
)
