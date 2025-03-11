package consumers

import (
	"GoLoad/internal/dataaccess/mq/producer"
	"GoLoad/internal/logic"
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
	downloadTaskLogic logic.DownloadTask
	logger            *zap.Logger
}

func NewDownloadTaskCreated(
	logger *zap.Logger,
	downloadTaskLogic logic.DownloadTask,
) DownloadTaskCreated {
	return &downloadTaskCreatedHandler{
		logger:            logger,
		downloadTaskLogic: downloadTaskLogic,
	}
}

func (d downloadTaskCreatedHandler) Handle(ctx context.Context, event producer.DownloadTaskCreated) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("event", event))
	logger.Info("Download task created event received")

	if err := d.downloadTaskLogic.ExecuteDownloadTask(ctx, event.ID); err != nil {
		logger.With(zap.Error(err)).Error("Failed to handle download task created event")
		return err
	}

	return nil
}
