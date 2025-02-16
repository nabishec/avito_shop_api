package send

import (
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/model"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type PostSendCoins interface {
	SendCoinsToUser(sendData model.SendCoinRequest, userID uuid.UUID) error
	UserNameExist(userName string) error
	UserIDExist(userID uuid.UUID) error
}
