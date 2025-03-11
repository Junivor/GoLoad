package logic

import (
	"GoLoad/internal/utils"
	"context"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	HTTPResponseHeaderContentType = "Content-Type"
	HTTPMetadataKeyContentType    = "content-type"
)

type Downloader interface {
	Download(ctx context.Context, writer io.Writer) (map[string]any, error)
}

type HTTPDownloader struct {
	url    string
	logger *zap.Logger
}

func NewHTTPDownloader(
	url string,
	logger *zap.Logger,
) Downloader {
	return &HTTPDownloader{
		url:    url,
		logger: logger,
	}
}

func (h HTTPDownloader) Download(ctx context.Context, writer io.Writer) (map[string]any, error) {
	logger := utils.LoggerWithContext(ctx, h.logger)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, http.NoBody)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to create http request")
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed make http request")
		return nil, err
	}

	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to read response and write to writer")
		return nil, err
	}

	metadata := map[string]any{
		HTTPMetadataKeyContentType: response.Header.Get(HTTPResponseHeaderContentType),
	}

	return metadata, nil
}
