package logic

import (
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/generated/grpc/go_load"
	"context"
	"database/sql"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateAccountParams struct {
	AccountName string
	Password    string
}

type CreateSessionParams struct {
	AccountName string
	Password    string
}

type CreateAccountOutput struct {
	ID          uint64
	AccountName string
}

type CreateSessionOutput struct {
	Account *go_load.Account
	Token   string
}

type Account interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
}

type account struct {
	goquDatabase                *goqu.Database
	accountDataAccessor         database.AccountDataAccessor
	accountDataPasswordAccessor database.AccountPasswordDataAccessor
	hashLogic                   Hash
	tokenLogic                  Token
}

func NewAccount(
	goquDatabase *goqu.Database,
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordAccessor database.AccountPasswordDataAccessor,
	hash Hash,
	tokenLogic Token,
) Account {
	return &account{
		goquDatabase:                goquDatabase,
		accountDataAccessor:         accountDataAccessor,
		accountDataPasswordAccessor: accountPasswordAccessor,
		hashLogic:                   hash,
		tokenLogic:                  tokenLogic,
	}
}

func (a account) databaseAccountToProtoAccount(account database.Account) *go_load.Account {
	return &go_load.Account{
		Id:          account.ID,
		AccountName: account.AccountName,
	}
}

func (a account) isAccountAccountNameTaken(ctx context.Context, accountName string) (bool, error) {
	if _, err := a.accountDataAccessor.GetAccountByAccountName(ctx, accountName); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a account) CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error) {
	var accountID uint64
	txErr := a.goquDatabase.WithTx(func(txDatabase *goqu.TxDatabase) error {
		accountNameTaken, err := a.isAccountAccountNameTaken(ctx, params.AccountName)

		if err != nil {
			return nil
		}

		if accountNameTaken {
			return errors.New("AccountName is already taken")
		}

		accountID, err := a.accountDataAccessor.WithDatabase(txDatabase).CreateAccount(ctx, database.Account{
			AccountName: params.AccountName,
		})

		if err != nil {
			return nil
		}

		hashedPassword, err := a.hashLogic.Hash(ctx, params.Password)

		if err != nil {
			return err
		}

		if err := a.accountDataPasswordAccessor.WithDatabase(txDatabase).CreateAccountPassword(ctx, database.AccountPassword{
			OfAccountID: accountID,
			Hash:        hashedPassword,
		}); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return CreateAccountOutput{}, txErr
	}

	return CreateAccountOutput{
		ID:          accountID,
		AccountName: params.AccountName,
	}, nil
}

func (a account) CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error) {
	existingAccount, err := a.accountDataAccessor.GetAccountByAccountName(ctx, params.AccountName)
	if err != nil {
		return CreateSessionOutput{}, err
	}

	existingAccountPassword, err := a.accountDataPasswordAccessor.GetAccountPassword(ctx, existingAccount.ID)
	if err != nil {
		return CreateSessionOutput{}, err
	}

	isHashEqual, err := a.hashLogic.IsHashEqual(ctx, params.Password, existingAccountPassword.Hash)
	if err != nil {
		return CreateSessionOutput{}, err
	}

	if !isHashEqual {
		return CreateSessionOutput{}, status.Error(codes.Unauthenticated, "incorrect password")
	}

	token, _, err := a.tokenLogic.GetToken(ctx, existingAccount.ID)
	if err != nil {
		return CreateSessionOutput{}, err
	}

	return CreateSessionOutput{
		Account: a.databaseAccountToProtoAccount(existingAccount),
		Token:   token,
	}, nil
}
