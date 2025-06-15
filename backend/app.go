package backend

import (
	"context"
	"maxxgui/backend/consts"
	"maxxgui/backend/handler"
	"maxxgui/backend/model"

	"github.com/google/wire"
)

var AppSet = wire.NewSet(wire.Struct(new(App), "*"))

type App struct {
	ctx          context.Context
	CrackHandler *handler.CrackHandler
}

func (a *App) Copyright() string {
	return consts.COPYWRIGHT
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	a.CrackHandler.Ctx = ctx
}

func (a *App) Bind() (binders []any) {
	binders = append(binders,
		a.CrackHandler,

		model.DefaultCrackProvider,
	)
	return
}

func (a *App) Enums() (enums []any) {
	enums = append(enums,
		consts.EventEnums,
	)
	return
}
