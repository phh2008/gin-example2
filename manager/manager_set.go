package manager

import "github.com/google/wire"

// ProviderSet is manager provider set
var ProviderSet = wire.NewSet(
	NewAuthorizeManager,
)
