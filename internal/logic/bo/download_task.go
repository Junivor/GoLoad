package bo

import "GoLoad/internal/generated/grpc/go_load"

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
