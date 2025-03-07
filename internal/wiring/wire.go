//go:build wireinject
// +build wireinject

// go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"GoLoad/internal/configs"
	"GoLoad/internal/dataaccess"
	"GoLoad/internal/handler"
	"GoLoad/internal/handler/grpc"
	"GoLoad/internal/logic"
	"GoLoad/internal/utils"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	dataaccess.WireSet,
	utils.WireSet,
	logic.WireSet,
	handler.WireSet,
)

func InitializeGRPCServer(path configs.ConfigFilePath) (grpc.Server, func(), error) {
	wire.Build(WireSet)
	return nil, nil, nil
}
