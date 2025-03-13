package models

import "github.com/doug-martin/goqu/v9"

var (
	TabNameTokenPublicKeys = goqu.T("token_public_keys")
)

const (
	ColNameTokenPublicKeysID        = "id"
	ColNameTokenPublicKeysPublicKey = "public_key"
)

type TokenPublicKey struct {
	ID        uint64 `db:"id" goqu:"skipinsert,skipupdate"`
	PublicKey string `db:"public_key"`
}
