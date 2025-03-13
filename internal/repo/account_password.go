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

type AccountPasswordDataAccessor interface {
	CreateAccountPassword(ctx context.Context, accountPassword models.AccountPassword) error
	GetAccountPassword(ctx context.Context, ofAccountID uint64) (models.AccountPassword, error)
	UpdateAccountPassword(ctx context.Context, accountPassword models.AccountPassword) error
	WithDatabase(database database.Database) AccountPasswordDataAccessor
}

type accountPasswordDataAccessor struct {
	database database.Database
	logger   *zap.Logger
}

func NewAccountPasswordDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a accountPasswordDataAccessor) CreateAccountPassword(ctx context.Context, accountPassword models.AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	_, err := a.database.
		Insert(models.TabNameAccountPasswords).
		Rows(goqu.Record{
			models.ColNameAccountPasswordsOfAccountID: accountPassword.OfAccountID,
			models.ColNameAccountPasswordsHash:        accountPassword.Hash,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create account password")
		return errors.ErrInternal("failed to create account password")
	}

	return nil
}

func (a accountPasswordDataAccessor) GetAccountPassword(
	ctx context.Context,
	ofAccountID uint64,
) (models.AccountPassword, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("of_account_id", ofAccountID))
	accountPassword := models.AccountPassword{}
	found, err := a.database.
		From(models.TabNameAccountPasswords).
		Where(goqu.Ex{models.ColNameAccountPasswordsOfAccountID: ofAccountID}).
		ScanStructContext(ctx, &accountPassword)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account password by id")
		return models.AccountPassword{}, errors.ErrInternal("failed to get account password by id")
	}

	if !found {
		logger.Warn("cannot find account by id")
		return models.AccountPassword{}, sql.ErrNoRows
	}

	return accountPassword, nil
}

func (a accountPasswordDataAccessor) UpdateAccountPassword(ctx context.Context, accountPassword models.AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	_, err := a.database.
		Update(models.TabNameAccountPasswords).
		Set(goqu.Record{models.ColNameAccountPasswordsHash: accountPassword.Hash}).
		Where(goqu.Ex{models.ColNameAccountPasswordsOfAccountID: accountPassword.OfAccountID}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update account password")
		return errors.ErrInternal("failed to update account password")
	}

	return nil
}

func (a accountPasswordDataAccessor) WithDatabase(database database.Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
		logger:   a.logger,
	}
}
