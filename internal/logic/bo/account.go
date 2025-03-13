package bo

import "GoLoad/internal/generated/grpc/go_load"

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
