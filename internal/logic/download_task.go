package logic

import (
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/dataaccess/mq/producer"
	"GoLoad/internal/generated/grpc/go_load"
	"context"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type CreateDownloadTaskParams struct {
	Token        string
	DownloadType go_load.DownloadType
	URL          string
}

type CreateDownloadTaskOutput struct {
	DownloadTask go_load.DownloadTask
}

type GetDownloadTaskListParams struct {
	Token  string
	offset uint64
	limit  uint64
}

type GetDownloadTaskListOutput struct {
	DownloadTask           go_load.DownloadTask
	TotalDownloadTaskCount uint64
}

type UpdateDownloadTaskParams struct {
	Token          string
	DownloadTaskID uint64
	url            string
}

type UpdateDownloadTaskOutput struct {
	DownloadTask go_load.DownloadTask
}

type DeleteDownloadTaskParams struct {
	Token          string
	DownloadTaskID uint64
}

type DeleteDownloadTaskOutput struct{}

type DownloadTask interface {
	CreateDownloadTask(ctx context.Context, params CreateDownloadTaskParams) (CreateDownloadTaskOutput, error)
	GetDownloadTaskList(ctx context.Context, params GetDownloadTaskListParams) GetDownloadTaskListOutput
	UpdateDownloadTask(ctx context.Context, params UpdateDownloadTaskParams) UpdateDownloadTaskOutput
	DeleteDownloadTask(ctx context.Context, params DeleteDownloadTaskParams) DeleteDownloadTaskOutput
}

type downLoadTask struct {
	logger                      *zap.Logger
	goquDatabase                *goqu.Database
	tokenLogic                  Token
	downloadTaskAccessor        database.DownloadTaskDataAccessor
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer
}

func NewDownloadTask(
	logger *zap.Logger,
	goquDatabase *goqu.Database,
	tokenLogic Token,
	downloadTaskAccessor database.DownloadTaskDataAccessor,
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer,
) DownloadTask {
	return &downLoadTask{
		logger:                      logger,
		goquDatabase:                goquDatabase,
		tokenLogic:                  tokenLogic,
		downloadTaskAccessor:        downloadTaskAccessor,
		downloadTaskCreatedProducer: downloadTaskCreatedProducer,
	}
}

func (d downLoadTask) CreateDownloadTask(ctx context.Context, params CreateDownloadTaskParams) (CreateDownloadTaskOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return CreateDownloadTaskOutput{}, err
	}

	databaseDownLoadTask := database.DownloadTask{
		OfAccountID:    accountID,
		DownloadType:   params.DownloadType,
		URL:            params.URL,
		DownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		Metadata:       "{}",
	}

	txErr := d.goquDatabase.WithTx(func(txDatabase *goqu.TxDatabase) error {
		downloadTaskID, createdDownloadTaskErr := d.downloadTaskAccessor.
			WithDatabase(txDatabase).
			CreateDownloadTask(ctx, databaseDownLoadTask)
		if createdDownloadTaskErr != nil {
			return createdDownloadTaskErr
		}

		databaseDownLoadTask.ID = downloadTaskID
		err = d.downloadTaskCreatedProducer.Produce(ctx, producer.DownloadTaskCreated{
			DownloadTask: databaseDownLoadTask,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return CreateDownloadTaskOutput{}, txErr
	}

	return CreateDownloadTaskOutput{
		DownloadTask: go_load.DownloadTask{
			Id:             databaseDownLoadTask.ID,
			OfAccount:      nil,
			DownloadType:   params.DownloadType,
			Url:            params.URL,
			DownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		},
	}, nil
}

func (d downLoadTask) GetDownloadTaskList(ctx context.Context, params GetDownloadTaskListParams) GetDownloadTaskListOutput {
	//TODO implement me
	panic("implement me")
}

func (d downLoadTask) UpdateDownloadTask(ctx context.Context, params UpdateDownloadTaskParams) UpdateDownloadTaskOutput {
	//TODO implement me
	panic("implement me")
}

func (d downLoadTask) DeleteDownloadTask(ctx context.Context, params DeleteDownloadTaskParams) DeleteDownloadTaskOutput {
	//TODO implement me
	panic("implement me")
}
