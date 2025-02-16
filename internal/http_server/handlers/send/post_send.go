package send

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/http_server/middlweare"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/nabishec/avito_shop_api/internal/storage/db"
	"github.com/rs/zerolog/log"
)

type SendingCoins struct {
	postSendCoins PostSendCoins
}

func NewSendingCoins(postSendCoins PostSendCoins) SendingCoins {
	return SendingCoins{
		postSendCoins: postSendCoins,
	}
}

// @Summary Отправить монеты другому пользователю.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.SendCoinRequest true "Данные для отправки монет"
// @Success 200 "Успешный ответ."
// @Failure 400 {object} model.ErrorResponse "Неверный запрос."
// @Failure 401 {object} model.ErrorResponse "Неавторизован."
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера."
// @Router /api/sendCoin [post]
func (h *SendingCoins) SendCoins(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http_server.hadnlers.send.BuyingItemByUser()"

	logger := log.With().Str("fn", op).Logger()
	logger.Debug().Msg("Request for send coins has been received")

	var sendData model.SendCoinRequest

	err := json.NewDecoder(r.Body).Decode(&sendData)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to decode request body")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}
	logger.Debug().Msg("Request body decoded")

	err = validator.New().Struct(sendData)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to validate request body")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}

	userIDStr := r.Context().Value(middlweare.RequestUserIDKey).(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error().Msg("Failed to get userID")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	if sendData.Amount <= 0 {
		log.Error().Msg("Amount < 0")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}

	err = h.postSendCoins.UserIDExist(userID)
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

	err = h.postSendCoins.UserNameExist(sendData.ToUser)
	if err != nil {
		if err == db.ErrUserNameNotExist {
			log.Error().Err(err)

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
			return
		}
		logger.Error().Err(err).Msg("Failed check user name in db")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	err = h.postSendCoins.SendCoinsToUser(sendData, userID)
	if err != nil {
		if err == db.ErrNotEnoughCoins {
			logger.Error().Err(err)

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
			return
		}
		logger.Error().Err(err).Msg("Failed to send coins")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	w.WriteHeader(http.StatusOK)
}
