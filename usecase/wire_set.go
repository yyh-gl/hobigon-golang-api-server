package usecase

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewNotificationUseCase,
	NewBirthdayUseCase,
	NewBlogUseCase,
)
