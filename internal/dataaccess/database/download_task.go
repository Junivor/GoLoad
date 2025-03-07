package database

import (
	"GoLoad/internal/generated/grpc/go_load"
	"context"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type DownloadTask struct {
	ID             uint64                 `sql:"id"`
	OfAccountID    uint64                 `sql:"of_account_id"`
	DownloadType   go_load.DownloadType   `sql:"download_type"`
	URL            string                 `sql:"url"`
	DownloadStatus go_load.DownloadStatus `sql:"download_status"`
	Metadata       string                 `sql:"metadata"`
}

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, task DownloadTask) (uint64, error)
	GetDownloadTaskList(ctx context.Context, userID, offset, limit uint64) ([]DownloadTask, error)
	GetDownloadTaskOfUser(ctx context.Context, userID uint64) (uint64, error)
	UpdateDownloadTask(ctx context.Context, task DownloadTask) error
	DeleteDownloadTask(ctx context.Context, id uint64) error
	WithDatabase(database Database) DownloadTaskDataAccessor
}

type downloadTaskDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func (d downloadTaskDataAccessor) WithDatabase(database Database) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
	}
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

func (d downloadTaskDataAccessor) CreateDownloadTask(ctx context.Context, task DownloadTask) (uint64, error) {
	return 1, nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskList(ctx context.Context, userID, offset, limit uint64) ([]DownloadTask, error) {
	//TODO implement me
	panic("implement me")
}

func (d downloadTaskDataAccessor) GetDownloadTaskOfUser(ctx context.Context, userID uint64) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (d downloadTaskDataAccessor) UpdateDownloadTask(ctx context.Context, task DownloadTask) error {
	//TODO implement me
	panic("implement me")
}

func (d downloadTaskDataAccessor) DeleteDownloadTask(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}
