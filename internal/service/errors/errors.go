package svcErr

import "errors"

var (
	ErrWalletNotFound = errors.New("wallet not found")

	ErrInvalidParams = errors.New("invalid parameters provided")
)
