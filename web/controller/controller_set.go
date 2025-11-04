package controller

import (
	"github.com/google/wire"
)

// ProviderSet is controller provider set
var ProviderSet = wire.NewSet(
	NewHelloController,
	NewPermissionController,
	NewRoleController,
	NewUserController,
)
