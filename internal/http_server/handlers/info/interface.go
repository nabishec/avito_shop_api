package info

import (
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/model"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
type GetInfo interface {
	GetUserInfo(userID uuid.UUID) (userInfo *model.InfoResponse, err error)
	UserIDExist(userID uuid.UUID) error
}
