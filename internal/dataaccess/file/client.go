package file

import (
	"GoLoad/internal/configs"
	"GoLoad/internal/utils"
	"bufio"
	"context"
	"github.com/minio/minio-go"
	_ "github.com/minio/minio-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path"
)

type Client interface {
	Write(ctx context.Context, filePath string) (io.WriteCloser, error)
	Read(ctx context.Context, filePath string) (io.ReadCloser, error)
}
type bufferedFileReader struct {
	file           *os.File
	bufferedReader io.Reader
}

func newBufferedFileReader(
	file *os.File,
) io.ReadCloser {
	return &bufferedFileReader{
		file:           file,
		bufferedReader: bufio.NewReader(file),
	}
}

func (b bufferedFileReader) Read(p []byte) (int, error) {
	return b.bufferedReader.Read(p)
}
func (b bufferedFileReader) Close() error {
	return b.file.Close()
}

type LocalClient struct {
	downloadDirectory string
	logger            *zap.Logger
}

type S3Client struct {
	minioClient *minio.Client
	bucket      string
	logger      *zap.Logger
}

type s3ClientReadWriteCloser struct {
	writtenData []byte
	isClose     bool
}

func NewLocalClient(
	downloadConfig configs.Download,
	logger *zap.Logger,
) Client {
	return &LocalClient{
		downloadDirectory: downloadConfig.DownloadDirectory,
		logger:            logger,
	}
}

func NewS3Client(
	downloadConfig configs.Download,
	logger *zap.Logger,
) (Client, error) {
	minioClient, err := minio.New(downloadConfig.Address, downloadConfig.Username, downloadConfig.Password, false)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to create minio client")
		return nil, err
	}
	return &S3Client{
		minioClient: minioClient,
		bucket:      downloadConfig.Bucket,
		logger:      logger,
	}, nil
}

func newS3ClientWriteCloser(
	minioClient *minio.Client,
	bucketName string,
	objectName string,
) io.ReadWriteCloser {

	readWriteCloser := &s3ClientReadWriteCloser{
		writtenData: make([]byte, 0),
		isClose:     false,
	}

	go func() {
		_, _ = minioClient.PutObject(bucketName, objectName, readWriteCloser, -1, minio.PutObjectOptions{})
	}()

	return readWriteCloser
}

func (s s3ClientReadWriteCloser) Close() error {
	s.isClose = true
	return nil
}

func (s s3ClientReadWriteCloser) Read(p []byte) (int, error) {
	if len(s.writtenData) > 0 {
		writtenLength := copy(p, s.writtenData)
		s.writtenData = s.writtenData[writtenLength:]
		return writtenLength, nil
	}

	if s.isClose {
		return 0, io.EOF
	}

	return 0, nil
}

func (s s3ClientReadWriteCloser) Write(p []byte) (int, error) {
	s.writtenData = append(s.writtenData, p...)

	return len(p), nil
}

func (l LocalClient) Read(ctx context.Context, filePath string) (io.ReadCloser, error) {
	logger := utils.LoggerWithContext(ctx, l.logger).With(zap.String("file_path", filePath))

	absolutePath := path.Join(l.downloadDirectory, filePath)
	file, err := os.Open(absolutePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to open file")
		return nil, status.Error(codes.Internal, "Failed to open file")
	}

	return newBufferedFileReader(file), nil
}
func (l LocalClient) Write(ctx context.Context, filePath string) (io.WriteCloser, error) {
	logger := utils.LoggerWithContext(ctx, l.logger).With(zap.String("file_path", filePath))

	absolutePath := path.Join(l.downloadDirectory, filePath)
	file, err := os.Create(absolutePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to open file")
		return nil, status.Error(codes.Internal, "Failed to open file")
	}

	return file, nil
}
func (s S3Client) Write(ctx context.Context, filePath string) (io.WriteCloser, error) {
	return newS3ClientWriteCloser(s.minioClient, s.bucket, filePath), nil
}
func (s S3Client) Read(ctx context.Context, filePath string) (io.ReadCloser, error) {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("file_path", filePath))
	object, err := s.minioClient.GetObject(s.bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get S3 object")
		return nil, status.Error(codes.Internal, "Failed to get S3 object")
	}

	return object, nil
}
