package logic

import (
	"GoLoad/internal/configs"
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/dataaccess/file"
	"GoLoad/internal/dataaccess/mq/producer"
	"GoLoad/internal/generated/grpc/go_load"
	"GoLoad/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/gammazero/workerpool"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

const (
	downloadTaskMetadataFieldNameFileName = "file-name"
)

type CreateDownloadTaskParams struct {
	Token        string
	DownloadType go_load.DownloadType
	URL          string
}

type CreateDownloadTaskOutput struct {
	DownloadTask *go_load.DownloadTask
}

type GetDownloadTaskListParams struct {
	Token  string
	Offset uint64
	Limit  uint64
}

type GetDownloadTaskListOutput struct {
	DownloadTaskList       []*go_load.DownloadTask
	TotalDownloadTaskCount uint64
}

type UpdateDownloadTaskParams struct {
	Token          string
	DownloadTaskID uint64
	URL            string
}

type UpdateDownloadTaskOutput struct {
	DownloadTask *go_load.DownloadTask
}

type DeleteDownloadTaskParams struct {
	Token          string
	DownloadTaskID uint64
}

type GetDownloadTaskFileParams struct {
	Token          string
	DownloadTaskID uint64
}

type DeleteDownloadTaskOutput struct{}

type DownloadTask interface {
	CreateDownloadTask(context.Context, CreateDownloadTaskParams) (CreateDownloadTaskOutput, error)
	GetDownloadTaskList(context.Context, GetDownloadTaskListParams) (GetDownloadTaskListOutput, error)
	UpdateDownloadTask(context.Context, UpdateDownloadTaskParams) (UpdateDownloadTaskOutput, error)
	DeleteDownloadTask(context.Context, DeleteDownloadTaskParams) error
	ExecuteAllPendingDownloadTask(context.Context) error
	ExecuteDownloadTask(context.Context, uint64) error
	GetDownloadTaskFile(context.Context, GetDownloadTaskFileParams) (io.ReadCloser, error)
	UpdateDownloadingAndFailedDownloadTaskStatusToPending(context.Context) error
}

type downLoadTask struct {
	logger                      *zap.Logger
	goquDatabase                *goqu.Database
	tokenLogic                  Token
	accountDataAccessor         database.AccountDataAccessor
	downloadTaskDataAccessor    database.DownloadTaskDataAccessor
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer
	fileClient                  file.Client
	cronConfig                  configs.Cron
}

func NewDownloadTask(
	logger *zap.Logger,
	goquDatabase *goqu.Database,
	tokenLogic Token,
	accountDataAccessor database.AccountDataAccessor,
	downloadTaskDataAccessor database.DownloadTaskDataAccessor,
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer,
	fileClient file.Client,
	cronConfig configs.Cron,
) DownloadTask {
	return &downLoadTask{
		logger:                      logger,
		goquDatabase:                goquDatabase,
		tokenLogic:                  tokenLogic,
		accountDataAccessor:         accountDataAccessor,
		downloadTaskDataAccessor:    downloadTaskDataAccessor,
		downloadTaskCreatedProducer: downloadTaskCreatedProducer,
		fileClient:                  fileClient,
		cronConfig:                  cronConfig,
	}
}

func (d downLoadTask) databaseDownloadTaskToProtoDownloadTask(
	downloadTask database.DownloadTask,
	account database.Account,
) *go_load.DownloadTask {
	return &go_load.DownloadTask{
		Id: downloadTask.ID,
		OfAccount: &go_load.Account{
			Id:          account.ID,
			AccountName: account.AccountName,
		},
		DownloadType:   downloadTask.DownloadType,
		Url:            downloadTask.URL,
		DownloadStatus: downloadTask.DownloadStatus,
	}
}

func (d downLoadTask) CreateDownloadTask(ctx context.Context, params CreateDownloadTaskParams) (CreateDownloadTaskOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return CreateDownloadTaskOutput{}, err
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		return CreateDownloadTaskOutput{}, err
	}

	downloadTask := database.DownloadTask{
		OfAccountID:    accountID,
		DownloadType:   params.DownloadType,
		URL:            params.URL,
		DownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		Metadata: database.JSON{
			Data: make(map[string]any),
		},
	}

	txErr := d.goquDatabase.WithTx(func(txDatabase *goqu.TxDatabase) error {
		downloadTaskID, createdDownloadTaskErr := d.downloadTaskDataAccessor.
			WithDatabase(txDatabase).
			CreateDownloadTask(ctx, downloadTask)
		if createdDownloadTaskErr != nil {
			return createdDownloadTaskErr
		}

		downloadTask.ID = downloadTaskID
		err = d.downloadTaskCreatedProducer.Produce(ctx, producer.DownloadTaskCreated{
			ID: downloadTask.ID,
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
		DownloadTask: d.databaseDownloadTaskToProtoDownloadTask(downloadTask, account),
	}, nil
}

func (d downLoadTask) GetDownloadTaskList(
	ctx context.Context,
	params GetDownloadTaskListParams,
) (GetDownloadTaskListOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	totalDownloadTaskCount, err := d.downloadTaskDataAccessor.GetDownloadTaskCountOfAccount(ctx, accountID)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	downloadTaskList, err := d.downloadTaskDataAccessor.
		GetDownloadTaskListOfAccount(ctx, accountID, params.Offset, params.Limit)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	return GetDownloadTaskListOutput{
		TotalDownloadTaskCount: totalDownloadTaskCount,
		DownloadTaskList: lo.Map(downloadTaskList, func(item database.DownloadTask, _ int) *go_load.DownloadTask {
			return d.databaseDownloadTaskToProtoDownloadTask(item, account)
		}),
	}, nil
}

func (d downLoadTask) UpdateDownloadTask(
	ctx context.Context,
	params UpdateDownloadTaskParams,
) (UpdateDownloadTaskOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return UpdateDownloadTaskOutput{}, err
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		return UpdateDownloadTaskOutput{}, err
	}

	output := UpdateDownloadTaskOutput{}
	txErr := d.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		downloadTask, getDownloadTaskWithXLockErr := d.downloadTaskDataAccessor.WithDatabase(td).
			GetDownloadTaskWithXLock(ctx, params.DownloadTaskID)
		if getDownloadTaskWithXLockErr != nil {
			return getDownloadTaskWithXLockErr
		}

		if downloadTask.OfAccountID != accountID {
			return status.Error(codes.PermissionDenied, "trying to update a download task the account does not own")
		}

		downloadTask.URL = params.URL
		output.DownloadTask = d.databaseDownloadTaskToProtoDownloadTask(downloadTask, account)
		return d.downloadTaskDataAccessor.WithDatabase(td).UpdateDownloadTask(ctx, downloadTask)
	})
	if txErr != nil {
		return UpdateDownloadTaskOutput{}, txErr
	}

	return output, nil
}

func (d downLoadTask) DeleteDownloadTask(ctx context.Context, params DeleteDownloadTaskParams) error {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}

	return d.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		downloadTask, getDownloadTaskWithXLockErr := d.downloadTaskDataAccessor.WithDatabase(td).
			GetDownloadTaskWithXLock(ctx, params.DownloadTaskID)
		if getDownloadTaskWithXLockErr != nil {
			return getDownloadTaskWithXLockErr
		}

		if downloadTask.OfAccountID != accountID {
			return status.Error(codes.PermissionDenied, "trying to delete a download task the account does not own")
		}

		return d.downloadTaskDataAccessor.WithDatabase(td).DeleteDownloadTask(ctx, params.DownloadTaskID)
	})
}

func (d downLoadTask) ExecuteAllPendingDownloadTask(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, d.logger)

	pendingDownloadTaskIDList, err := d.downloadTaskDataAccessor.GetPendingDownloadTaskIDList(ctx)
	if err != nil {
		return err
	}
	if len(pendingDownloadTaskIDList) == 0 {
		logger.Info("no pending download task found")
		return nil
	}

	logger.
		With(zap.Int("len(pending_download_task_id_list)", len(pendingDownloadTaskIDList))).
		Info("pending download task found")

	workerPool := workerpool.New(d.cronConfig.ExecuteAllPendingDownloadTask.ConcurrencyLimit)
	for _, id := range pendingDownloadTaskIDList {
		workerPool.Submit(func() {
			if executeDownloadTaskErr := d.ExecuteDownloadTask(ctx, id); executeDownloadTaskErr != nil {
				logger.
					With(zap.Uint64("download_task_id", id)).
					With(zap.Error(executeDownloadTaskErr)).
					Error("failed to execute download_task")
			}
		})
	}

	workerPool.StopWait()
	return nil
}

func (d downLoadTask) updateDownloadTaskFromPendingToDownloading(ctx context.Context, id uint64) (
	bool,
	database.DownloadTask,
	error,
) {

	var (
		logger       = utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))
		updated      = false
		downLoadTask database.DownloadTask
		err          error
	)

	if txErr := d.goquDatabase.WithTx(func(txDatabase *goqu.TxDatabase) error {
		downloadTask, err := d.downloadTaskDataAccessor.WithDatabase(txDatabase).GetDownloadTaskWithXLock(ctx, id)
		if err != nil {
			if errors.Is(err, database.ErrDownloadTaskNotFound) {
				logger.Warn("Download task not found, will skip")
				return nil
			}

			logger.With(zap.Error(err)).Error("Failed to get download task")
			return err
		}

		if downloadTask.DownloadStatus != go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING {
			logger.Warn("Download task is not in pending status, will not execute")
			updated = false
			return nil
		}

		downloadTask.DownloadStatus = go_load.DownloadStatus_DOWNLOAD_STATUS_DOWNLOADING
		if err := d.downloadTaskDataAccessor.WithDatabase(txDatabase).UpdateDownloadTask(ctx, downloadTask); err != nil {
			logger.With(zap.Error(err)).Error("Failed to update download task")
			return err
		}

		updated = true
		return nil
	}); txErr != nil {
		return false, database.DownloadTask{}, err
	}

	return updated, downLoadTask, nil
}

func (d downLoadTask) ExecuteDownloadTask(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	updated, downloadTask, err := d.updateDownloadTaskFromPendingToDownloading(ctx, id)
	if err != nil {
		return err
	}
	if !updated {
		return nil
	}

	var downloader Downloader

	switch downloadTask.DownloadType {
	case go_load.DownloadType_DOWNLOAD_TYPE_HTTP:
		downloader = NewHTTPDownloader(downloadTask.URL, d.logger)
	default:
		logger.With(zap.Any("download_type", downloadTask.DownloadType)).Error("Unsupported download type")
	}

	fileName := fmt.Sprintf("download_file_%d", id)
	fileWriteCloser, err := d.fileClient.Write(ctx, fileName)

	if err != nil {
		return err
	}

	defer fileWriteCloser.Close()

	metadata, err := downloader.Download(ctx, nil)

	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to download")
		return err
	}

	metadata["file-name"] = fileName
	downloadTask.DownloadStatus = go_load.DownloadStatus_DOWNLOAD_STATUS_SUCCESS
	downloadTask.Metadata = database.JSON{Data: metadata}

	if err := d.downloadTaskDataAccessor.UpdateDownloadTask(ctx, downloadTask); err != nil {
		logger.With(zap.Error(err)).Error("Failed to update download task status to success")
		return err
	}

	logger.Info("Download task executed successfully")

	return nil
}

func (d downLoadTask) GetDownloadTaskFile(
	ctx context.Context,
	params GetDownloadTaskFileParams,
) (io.ReadCloser, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return nil, err
	}

	downloadTask, err := d.downloadTaskDataAccessor.GetDownloadTask(ctx, params.DownloadTaskID)
	if err != nil {
		return nil, err
	}

	if downloadTask.OfAccountID != accountID {
		return nil, status.Error(codes.PermissionDenied, "trying to get file of a download task the account does not own")
	}

	if downloadTask.DownloadStatus != go_load.DownloadStatus_DOWNLOAD_STATUS_SUCCESS {
		return nil, status.Error(codes.InvalidArgument, "download task does not have status of success")
	}

	downloadTaskMetadata, ok := downloadTask.Metadata.Data.(map[string]any)
	if !ok {
		return nil, status.Error(codes.Internal, "download task metadata is not a map[string]any")
	}

	fileName, ok := downloadTaskMetadata[downloadTaskMetadataFieldNameFileName]
	if !ok {
		return nil, status.Error(codes.Internal, "download task metadata does not contain file name")
	}

	return d.fileClient.Read(ctx, fileName.(string))
}

func (d downLoadTask) UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error {
	return d.downloadTaskDataAccessor.UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx)
}
