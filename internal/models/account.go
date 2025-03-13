package models

import (
	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameAccounts = goqu.T("accounts")
)

const (
	ColNameAccountsID          = "id"
	ColNameAccountsAccountName = "account_name"
)

type Account struct {
	ID          uint64 `db:"id" goqu:"skipinsert,skipupdate"`
	AccountName string `db:"account_name"`
}
