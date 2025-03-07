package consumer

import (
	"GoLoad/internal/configs"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

type HandlerFunc func(ctx context.Context, queueName string, payload []byte) error

type Consumer interface {
	RegisterHandler(queueName string, handlerFunc HandlerFunc) error
	Start(ctx context.Context) error
}

type partitionConsumerAndHandlerFunc struct {
	queueName         string
	partitionConsumer sarama.PartitionConsumer
	handlerFunc       HandlerFunc
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = mqConfig.ClientID
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

type consumer struct {
	saramaConsumer                      sarama.Consumer
	partitionConsumerAndHandlerFuncList []partitionConsumerAndHandlerFunc
	logger                              *zap.Logger
}

func NewConsumer(
	mqConfig configs.MQ,
	logger *zap.Logger,
) (Consumer, error) {
	saramaConsumer, err := sarama.NewConsumer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		return nil, fmt.Errorf("Failed to create sarama consumers: %w", err)
	}

	return &consumer{
		saramaConsumer: saramaConsumer,
		logger:         logger,
	}, nil
}

func (c consumer) RegisterHandler(queueName string, handlerFunc HandlerFunc) error {
	partitionConsumer, err := c.saramaConsumer.ConsumePartition(queueName, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("Failed to create sarama partition consumers: %w", err)
	}

	c.partitionConsumerAndHandlerFuncList = append(
		c.partitionConsumerAndHandlerFuncList,
		partitionConsumerAndHandlerFunc{
			queueName:         queueName,
			partitionConsumer: partitionConsumer,
			handlerFunc:       handlerFunc,
		},
	)

	return nil
}

func (c consumer) Start(_ context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _, partitionConsumerAndHandler := range c.partitionConsumerAndHandlerFuncList {
		go func() {
			queueName := partitionConsumerAndHandler.queueName
			partitionConsumer := partitionConsumerAndHandler.partitionConsumer
			handlerFunc := partitionConsumerAndHandler.handlerFunc
			logger := c.logger.With(zap.String("queue_name", queueName))

			for {
				select {
				case message := <-partitionConsumer.Messages():
					if err := handlerFunc(context.Background(), queueName, message.Value); err != nil {
						logger.With(zap.Error(err)).Error("Failed to handle message")
					}

				case <-signals:
					return
				default:

				}
			}
		}()
	}

	<-signals
	return nil
}
