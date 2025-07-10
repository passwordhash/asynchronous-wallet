package repoErr

import "errors"

var (
	ErrWalletNotFound = errors.New("wallet not found")
)
