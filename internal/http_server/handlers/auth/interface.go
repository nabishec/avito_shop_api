package auth

import (
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/model"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type PostAuth interface {
	GetUserID(model.AuthRequest) (uuid.UUID, error)
}
