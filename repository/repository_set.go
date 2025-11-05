package repository

import (
	"github.com/google/wire"
)

// ProviderSet is repository provider set
var ProviderSet = wire.NewSet(
	NewPermissionRepository,
	NewRoleRepository,
	NewRolePermissionRepository,
	NewUserRepository,
	NewUserRoleRepository,
)
