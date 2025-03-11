package dataaccess

import (
	"GoLoad/internal/dataaccess/cache"
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/dataaccess/file"
	"GoLoad/internal/dataaccess/mq"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	database.WireSet,
	cache.WireSet,
	mq.WireSet,
	file.WireSet,
)
