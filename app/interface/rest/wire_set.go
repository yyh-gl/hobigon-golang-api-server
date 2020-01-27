package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

// WireSet : interface層のWireSet（API用）
var WireSet = wire.NewSet(
	validator.New,
	NewNotification,
	NewBirthday,
	NewBlog,
)
