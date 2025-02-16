package buy

import (
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type GetBuy interface {
	GetItemByUser(userID uuid.UUID, item string) error
	UserIDExist(userID uuid.UUID) error
}
