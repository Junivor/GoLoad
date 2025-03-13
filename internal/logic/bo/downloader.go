package bo

import "go.uber.org/zap"

type HTTPDownloader struct {
	url    string
	logger *zap.Logger
}
