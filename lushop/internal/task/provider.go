package task

import "github.com/google/wire"

// ProviderSet provides asynq component.
var ProviderSet = wire.NewSet(NewAsynqComponent)
