package models

import (
	"GoLoad/internal/generated/grpc/go_load"
	"github.com/doug-martin/goqu/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	TabNameDownloadTasks = goqu.T("download_tasks")

	ErrDownloadTaskNotFound = status.Error(codes.NotFound, "download task not found")
)

const (
	ColNameDownloadTaskID             = "id"
	ColNameDownloadTaskOfAccountID    = "of_account_id"
	ColNameDownloadTaskDownloadType   = "download_type"
	ColNameDownloadTaskURL            = "url"
	ColNameDownloadTaskDownloadStatus = "download_status"
	ColNameDownloadTaskMetadata       = "metadata"
)

type DownloadTask struct {
	ID             uint64                 `db:"id" goqu:"skipinsert,skipupdate"`
	OfAccountID    uint64                 `db:"of_account_id" goqu:"skipupdate"`
	DownloadType   go_load.DownloadType   `db:"download_type"`
	URL            string                 `db:"url"`
	DownloadStatus go_load.DownloadStatus `db:"download_status"`
	Metadata       JSON                   `db:"metadata"`
}
