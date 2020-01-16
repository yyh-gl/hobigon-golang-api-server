package usecase

import "github.com/google/wire"

var APISet = wire.NewSet(
	NewNotificationUseCase,
	NewBirthdayUseCase,
	NewBlogUseCase,
)

var CLISet = wire.NewSet(
	NewNotificationUseCase,
)
