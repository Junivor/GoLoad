package repo

import (
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/errors"
	"GoLoad/internal/models"
	"GoLoad/internal/utils"
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account models.Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (models.Account, error)
	GetAccountByAccountName(ctx context.Context, accountName string) (models.Account, error)
	WithDatabase(database database.Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database database.Database
	logger   *zap.Logger
}

func NewAccountDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account models.Account) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("account", account))

	result, err := a.database.
		Insert(models.TabNameAccounts).
		Rows(goqu.Record{
			models.ColNameAccountsAccountName: account.AccountName,
		}).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create account")
		return 0, errors.ErrInternal("failed to create account")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, errors.ErrInternal("failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (models.Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	account := models.Account{}
	found, err := a.database.
		From(models.TabNameAccounts).
		Where(goqu.C(models.ColNameAccountsID).Eq(id)).
		ScanStructContext(ctx, &account)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account by id")
		return models.Account{}, errors.ErrNotFound("account")
	}

	if !found {
		logger.Warn("cannot find account by id")
		return models.Account{}, errors.ErrNotFound("account")
	}

	return account, nil
}

func (a accountDataAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (models.Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", accountName))

	account := models.Account{}
	found, err := a.database.
		From(models.TabNameAccounts).
		Where(goqu.C(models.ColNameAccountsAccountName).Eq(accountName)).
		ScanStructContext(ctx, &account)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account by name")
		return models.Account{}, errors.ErrInternal("failed to get account by name")
	}

	if !found {
		logger.Warn("cannot find account by name")
		return models.Account{}, errors.ErrNotFound("account")
	}

	return account, nil
}

func (a accountDataAccessor) WithDatabase(database database.Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   a.logger,
	}
}
