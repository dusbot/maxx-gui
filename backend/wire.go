//go:build wireinject

package backend

import (
	"context"
	"maxxgui/backend/handler"
	"maxxgui/backend/init_"

	"github.com/google/wire"
)

func NewApp() *App {
	wire.Build(
		AppSet,
		ProvideContext,
		init_.InitQuery,
		handler.ProviderSet,
	)
	return new(App)
}

func ProvideContext() context.Context {
	return context.TODO()
}
