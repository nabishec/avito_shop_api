package buy

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/nabishec/avito_shop_api/internal/http_server/middlweare"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/nabishec/avito_shop_api/internal/storage/db"
)

type Buying struct {
	getBuy GetBuy
}

func NewBuying(getBuy GetBuy) *Buying {
	return &Buying{
		getBuy: getBuy,
	}
}

// @Summary Купить предмет за монеты.
// @Security BearerAuth
// @Produce json
// @Success 200 "Успешный ответ."
// @Failure 400 {object} model.ErrorResponse "Неверный запрос."
// @Failure 401 {object} model.ErrorResponse "Неавторизован."
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера."
// @Router /api/buy/{item} [get]
func (h *Buying) BuyingItemByUser(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http_server.hadnlers.buy.BuyingItemByUser()"

	logger := log.With().Str("fn", op).Logger()
	logger.Debug().Msg("Request for buying has been received")

	userIDStr := r.Context().Value(middlweare.RequestUserIDKey).(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error().Msg("Failed to get userID")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	err = h.getBuy.UserIDExist(userID)
	if err != nil {
		if err == db.ErrUserIDNotExist {
			log.Error().Err(err)

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
			return
		}
		logger.Error().Err(err).Msg("Failed check user ID in db")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	item := chi.URLParam(r, "item")
	if item == "" {
		logger.Error().Msg("Failed to get parameters")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}
	logger.Debug().Msg("Parameters are received")

	err = h.getBuy.GetItemByUser(userID, item)
	if err != nil {
		if err == db.ErrNotEnoughCoins || err == db.ErrItemNotExist {
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
