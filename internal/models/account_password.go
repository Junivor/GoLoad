package models

import "github.com/doug-martin/goqu/v9"

var (
	TabNameAccountPasswords = goqu.T("account_passwords")
)

const (
	ColNameAccountPasswordsOfAccountID = "of_account_id"
	ColNameAccountPasswordsHash        = "hash"
)

type AccountPassword struct {
	OfAccountID uint64 `db:"of_account_id" goqu:"skipupdate"`
	Hash        string `db:"hash"`
}
