package repo

import (
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/errors"
	"GoLoad/internal/models"
	"GoLoad/internal/utils"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type TokenPublicKeyDataAccessor interface {
	CreatePublicKey(ctx context.Context, tokenPublicKey models.TokenPublicKey) (uint64, error)
	GetPublicKey(ctx context.Context, id uint64) (models.TokenPublicKey, error)
	WithDatabase(database database.Database) TokenPublicKeyDataAccessor
}

type tokenPublicKeyDataAccessor struct {
	database database.Database
	logger   *zap.Logger
}

func NewTokenPublicKeyDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a tokenPublicKeyDataAccessor) CreatePublicKey(
	ctx context.Context,
	tokenPublicKey models.TokenPublicKey,
) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	logger.Info("Inserting public key", zap.String("publicKey", tokenPublicKey.PublicKey))

	sqlQuery, args, _ := a.database.
		Insert(models.TabNameTokenPublicKeys).
		Rows(goqu.Record{
			models.ColNameTokenPublicKeysPublicKey: tokenPublicKey.PublicKey,
		}).
		ToSQL()
	logger.Info("Generated SQL Query", zap.String("query", sqlQuery), zap.Any("args", args))

	result, err := a.database.
		Insert(models.TabNameTokenPublicKeys).
		Rows(goqu.Record{
			models.ColNameTokenPublicKeysPublicKey: tokenPublicKey.PublicKey,
		}).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create token public key")
		return 0, errors.ErrInternal("failed to create token public key")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, errors.ErrInternal("failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (a tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, id uint64) (models.TokenPublicKey, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("id", id))

	tokenPublicKey := models.TokenPublicKey{}
	found, err := a.database.
		Select().
		From(models.TabNameTokenPublicKeys).
		Where(goqu.Ex{
			models.ColNameTokenPublicKeysID: id,
		}).
		Executor().
		ScanStructContext(ctx, &tokenPublicKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get public key")
		return models.TokenPublicKey{}, errors.ErrInternal("failed to get public key")
	}

	if !found {
		logger.Warn("public key not found")
		return models.TokenPublicKey{}, sql.ErrNoRows
	}

	return tokenPublicKey, nil
}

func (a tokenPublicKeyDataAccessor) WithDatabase(database database.Database) TokenPublicKeyDataAccessor {
	a.database = database
	return a
}
