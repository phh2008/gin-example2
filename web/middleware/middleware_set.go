package middleware

import "github.com/google/wire"

// ProviderSet is middleware provider set
var ProviderSet = wire.NewSet(NewAuth)
