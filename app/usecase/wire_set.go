package usecase

import "github.com/google/wire"

// APISet : usecase層のWireSet（API用）
var APISet = wire.NewSet(
	NewNotificationUseCase,
	NewBirthdayUseCase,
	NewBlogUseCase,
)

// CLISet : usecase層のWireSet（CLI用）
var CLISet = wire.NewSet(
	NewNotificationUseCase,
)
