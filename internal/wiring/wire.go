//go:build wireinject
// +build wireinject

//
//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"GoLoad/internal/app"
	"GoLoad/internal/configs"
	"GoLoad/internal/dataaccess"
	"GoLoad/internal/handler"
	"GoLoad/internal/logic"
	"GoLoad/internal/repo"
	"GoLoad/internal/utils"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	utils.WireSet,
	dataaccess.WireSet,
	logic.WireSet,
	handler.WireSet,
	app.WireSet,
	repo.WireSet,
)

func InitializeStandaloneServer(configFilePath configs.ConfigFilePath) (*app.StandaloneServer, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}
