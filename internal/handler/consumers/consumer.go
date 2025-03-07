package consumers

import (
	"GoLoad/internal/dataaccess/mq/consumer"
	"GoLoad/internal/dataaccess/mq/producer"
	"GoLoad/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type Root interface {
	Start(ctx context.Context) error
}

type root struct {
	downloadTaskCreatedHandler DownloadTaskCreated
	mqConsumer                 consumer.Consumer
	logger                     *zap.Logger
}

func NewRoot(
	downloadTaskCreated DownloadTaskCreated,
	mqConsumer consumer.Consumer,
	logger *zap.Logger,
) Root {
	return &root{
		downloadTaskCreatedHandler: downloadTaskCreated,
		mqConsumer:                 mqConsumer,
		logger:                     logger,
	}
}

func (r root) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, r.logger)
	if err := r.mqConsumer.RegisterHandler(
		producer.MessageQueueDownloadTaskCreated,
		func(ctx context.Context, queueName string, payload []byte) error {
			var event producer.DownloadTaskCreated
			if err := json.Unmarshal(payload, &event); err != nil {
				return err
			}

			return r.downloadTaskCreatedHandler.Handle(ctx, event)
		}); err != nil {
		logger.With(zap.Error(err)).Error("Failed to register message queue download")
		return fmt.Errorf("Failed to download task created handler: %w", err)
	}

	return r.mqConsumer.Start(ctx)
}
