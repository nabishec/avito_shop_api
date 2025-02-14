package db

import "errors"

var (
	ErrUserNotExist   = errors.New("user doesn't exist")
	ErrNotEnoughCoins = errors.New("user doesn't have enough coins to buy")
)
