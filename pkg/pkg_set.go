package pkg

import (
	"com.example/example/pkg/orm"
	"com.example/example/pkg/xcasbin"
	"com.example/example/pkg/xjwt"
	"github.com/google/wire"
)

// ProviderSet is pkg provider set
var ProviderSet = wire.NewSet(
	orm.NewDB,
	xjwt.NewJwtHelper,
	xcasbin.NewCasbin,
)
