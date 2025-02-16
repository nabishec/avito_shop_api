package info

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/nabishec/avito_shop_api/internal/http_server/middlweare"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/nabishec/avito_shop_api/internal/storage/db"
	"github.com/rs/zerolog/log"
)

type UserInformation struct {
	getInfo GetInfo
}

func NewUserInformation(getInfo GetInfo) UserInformation {
	return UserInformation{
		getInfo: getInfo,
	}
}

// @Summary Получить информацию о монетах, инвентаре и истории транзакций.
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.InfoResponse "Успешный ответ."
// @Failure 400 {object} model.ErrorResponse "Неверный запрос."
// @Failure 401 {object} model.ErrorResponse "Неавторизован."
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера."
// @Router /api/info [get]
func (h *UserInformation) ReturnUserInfo(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http_server.hadnlers.info.ReturnUserInfo()"

	logger := log.With().Str("fn", op).Logger()
	logger.Debug().Msg("Request for user's information has been received")

	userIDStr := r.Context().Value(middlweare.RequestUserIDKey).(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error().Msg("Failed to get userID")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	err = h.getInfo.UserIDExist(userID)
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

	response, err := h.getInfo.GetUserInfo(userID)
	if err != nil {
		logger.Error().Msg("Failed to get user information")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)

}
