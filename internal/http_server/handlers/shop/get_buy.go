package shop

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/nabishec/avito_shop_api/internal/storage/db"
	"github.com/rs/zerolog/log"
)

type GetBuy interface {
	GetItemByUser(userID uuid.UUID, item string) error
}

type Buying struct {
	getBuy GetBuy
}

func NewBuyItem(getBuy GetBuy) Buying {
	return Buying{
		getBuy: getBuy,
	}
}

func (h *Buying) BuyingItemByUser(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http_server.hadnlers.shop.BuyingItemByUser()"

	logger := log.With().Str("fn", op).Logger()
	logger.Debug().Msg("Request for buying has been received")

	item := r.URL.Query().Get("item")
	userID, err := uuid.FromBytes([]byte(r.URL.Query().Get("user_id")))

	if err != nil {
		logger.Error().Msg("Failed to get userID")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	if item == "" {
		logger.Error().Msg("Failed to get parameters")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}
	logger.Debug().Msg("Parameters are received")

	err = h.getBuy.GetItemByUser(userID, item)
	if err != nil {
		if err == db.ErrNotEnoughCoins {
			logger.Error().Err(err)

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
			return
		}
		logger.Error().Err(err).Msg("Failed to buy item")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	w.WriteHeader(http.StatusOK)
}
