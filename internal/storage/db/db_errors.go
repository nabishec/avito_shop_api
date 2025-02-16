package db

import "errors"

var (
	ErrNotEnoughCoins   = errors.New("user doesn't have enough coins to buy")
	ErrUserNameNotExist = errors.New("username not exist")
	ErrUserIDNotExist   = errors.New("user id not exist")
)
