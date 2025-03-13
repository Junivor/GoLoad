package repo

import (
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/errors"
	"GoLoad/internal/generated/grpc/go_load"
	"GoLoad/internal/models"
	"GoLoad/internal/utils"
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, task models.DownloadTask) (uint64, error)
	GetDownloadTaskListOfAccount(ctx context.Context, accountID, offset, limit uint64) ([]models.DownloadTask, error)
	GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error)
	GetDownloadTask(ctx context.Context, id uint64) (models.DownloadTask, error)
	GetDownloadTaskWithXLock(ctx context.Context, id uint64) (models.DownloadTask, error)
	UpdateDownloadTask(ctx context.Context, task models.DownloadTask) error
	DeleteDownloadTask(ctx context.Context, id uint64) error
	GetPendingDownloadTaskIDList(ctx context.Context) ([]uint64, error)
	UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error
	WithDatabase(database database.Database) DownloadTaskDataAccessor
}

type downloadTaskDataAccessor struct {
	database database.Database
	logger   *zap.Logger
}

func NewDownloadTaskDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (d downloadTaskDataAccessor) CreateDownloadTask(ctx context.Context, task models.DownloadTask) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("task", task))

	result, err := d.database.
		Insert(models.TabNameDownloadTasks).
		Rows(task).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create download task")
		return 0, status.Error(codes.Internal, "failed to create download task")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, status.Error(codes.Internal, "failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (d downloadTaskDataAccessor) DeleteDownloadTask(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	if _, err := d.database.
		Delete(models.TabNameDownloadTasks).
		Where(goqu.Ex{models.ColNameDownloadTaskID: id}).
		Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to delete download task")
		return status.Error(codes.Internal, "failed to delete download task")
	}

	return nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("account_id", accountID))

	count, err := d.database.
		From(models.TabNameDownloadTasks).
		Where(goqu.Ex{models.ColNameDownloadTaskOfAccountID: accountID}).
		CountContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to count download task of user")
		return 0, status.Error(codes.Internal, "failed to count download task of user")
	}

	return uint64(count), nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskListOfAccount(
	ctx context.Context,
	accountID uint64,
	offset uint64,
	limit uint64,
) ([]models.DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).
		With(zap.Uint64("account_id", accountID)).
		With(zap.Uint64("offset", offset)).
		With(zap.Uint64("limit", limit))

	downloadTaskList := make([]models.DownloadTask, 0)
	if err := d.database.
		Select().
		From(models.TabNameDownloadTasks).
		Where(goqu.Ex{models.ColNameAccountPasswordsOfAccountID: accountID}).
		Offset(uint(offset)).
		Limit(uint(limit)).
		Executor().
		ScanStructsContext(ctx, &downloadTaskList); err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task list of account")
		return nil, status.Error(codes.Internal, "failed to get download task list of account")
	}

	return downloadTaskList, nil
}

func (d downloadTaskDataAccessor) GetDownloadTask(ctx context.Context, id uint64) (models.DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	downloadTask := models.DownloadTask{}
	found, err := d.database.
		Select().
		From(models.TabNameDownloadTasks).
		Where(goqu.Ex{models.ColNameDownloadTaskID: id}).
		ScanStructContext(ctx, &downloadTask)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task")
		return models.DownloadTask{}, status.Error(codes.Internal, "failed to get download task list of account")
	}

	if !found {
		logger.Error("download task not found")
		return models.DownloadTask{}, errors.ErrNotFound("download task")
	}

	return downloadTask, nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskWithXLock(ctx context.Context, id uint64) (models.DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	downloadTask := models.DownloadTask{}
	found, err := d.database.
		Select().
		From(models.TabNameDownloadTasks).
		Where(goqu.Ex{models.ColNameDownloadTaskID: id}).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &downloadTask)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task")
		return models.DownloadTask{}, errors.ErrInternal("failed to get download task list of account")
	}

	if !found {
		logger.Error("download task not found")
		return models.DownloadTask{}, errors.ErrNotFound("download task")
	}

	return downloadTask, nil
}

func (d downloadTaskDataAccessor) UpdateDownloadTask(ctx context.Context, task models.DownloadTask) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("task", task))

	if _, err := d.database.
		Update(models.TabNameDownloadTasks).
		Set(task).
		Where(goqu.Ex{models.ColNameDownloadTaskID: task.ID}).
		Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to update download task")
		return errors.ErrInternal("failed to update download task")
	}

	return nil
}

func (d downloadTaskDataAccessor) GetPendingDownloadTaskIDList(ctx context.Context) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger)

	downloadTaskIDList := make([]uint64, 0)
	if err := d.database.
		Select(models.ColNameDownloadTaskID).
		From(models.TabNameDownloadTasks).
		Where(goqu.Ex{
			models.ColNameDownloadTaskDownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		}).
		ScanValsContext(ctx, &downloadTaskIDList); err != nil {
		logger.With(zap.Error(err)).Error("failed to get pending download task id list")
		return nil, errors.ErrInternal("failed to get pending download task id list")
	}

	return downloadTaskIDList, nil
}

func (d downloadTaskDataAccessor) UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, d.logger)

	if _, err := d.database.
		Update(models.TabNameDownloadTasks).
		Set(goqu.Record{
			models.ColNameDownloadTaskDownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		}).
		Where(
			goqu.C(models.ColNameDownloadTaskDownloadStatus).
				In(go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING, go_load.DownloadStatus_DOWNLOAD_STATUS_FAILED),
		).Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to update downloading and failed download task status to pending")
		return errors.ErrInternal("failed to update downloading and failed download task status to pending")
	}

	return nil
}

func (d downloadTaskDataAccessor) WithDatabase(database database.Database) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   d.logger,
	}
}
