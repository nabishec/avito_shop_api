package info

import (
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/model"
)

type GetInfo interface {
	GetUserInfo(userID uuid.UUID) (userInfo *model.InfoResponse, err error)
	UserIDExist(userID uuid.UUID) error
}
