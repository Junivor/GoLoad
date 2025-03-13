package bo

import (
	"GoLoad/internal/configs"
	"context"
)

type Hash interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, data string, hashed string) (bool, error)
}

type hash struct {
	authConfig configs.Auth
}
