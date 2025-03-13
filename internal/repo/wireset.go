package repo

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAccountDataAccessor,
	NewAccountPasswordDataAccessor,
	NewDownloadTaskDataAccessor,
	NewTokenPublicKeyDataAccessor,
)
