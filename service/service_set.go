package service

import (
	"github.com/google/wire"
)

// ProviderSet is service provider set
var ProviderSet = wire.NewSet(
	NewPermissionService,
	NewRolePermissionService,
	NewRoleService,
	NewUserService,
)
