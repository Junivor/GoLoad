package consumers

import (
	"GoLoad/internal/dataaccess/mq/producer"
	"GoLoad/internal/utils"
	"context"
	"go.uber.org/zap"
)

const (
	DownloadTaskCreatedMessageQueue = "download_task_created"
)

type DownloadTaskCreated interface {
	Handle(ctx context.Context, event producer.DownloadTaskCreated) error
}

type downloadTaskCreatedHandler struct {
	logger *zap.Logger
}

func NewDownloadTaskCreated(logger *zap.Logger) DownloadTaskCreated {
	return &downloadTaskCreatedHandler{logger: logger}
}

func (d downloadTaskCreatedHandler) Handle(ctx context.Context, event producer.DownloadTaskCreated) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("event", event))
	logger.Info("Download task created event received")

	return nil
}
