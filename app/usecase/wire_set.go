package usecase

import "github.com/google/wire"

// APISet : usecase層のWireSet（API用）
var APISet = wire.NewSet(
	NewNotification,
	NewBirthday,
	NewBlog,
)

// CLISet : usecase層のWireSet（CLI用）
var CLISet = wire.NewSet(
	NewNotification,
)
