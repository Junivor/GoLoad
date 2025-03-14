package handler

import (
	"GoLoad/internal/handler/consumers"
	"GoLoad/internal/handler/grpc"
	"GoLoad/internal/handler/http"
	"GoLoad/internal/handler/jobs"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	grpc.WireSet,
	http.WireSet,
	consumers.WireSet,
	jobs.WireSet,
)
