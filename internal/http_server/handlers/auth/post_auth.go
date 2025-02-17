package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/nabishec/avito_shop_api/internal/pkg/token"
	"github.com/nabishec/avito_shop_api/internal/storage/db"
	"github.com/rs/zerolog/log"
)

type Auth struct {
	postAuth PostAuth
}

func NewAuth(postAuth PostAuth) Auth {
	return Auth{
		postAuth: postAuth,
	}
}

// @Summary Аутентификация и получение JWT-токена.
// @Accept json
// @Produce json
// @Param body body model.AuthRequest true "Данные для аутентификации"
// @Success 200 {object} model.AuthResponse "Успешная аутентификация."
// @Failure 400 {object} model.ErrorResponse "Неверный запрос."
// @Failure 401 {object} model.ErrorResponse "Неавторизован."
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера."
// @Router /api/auth [post]
func (h *Auth) ReturnAuthToken(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http_server.hadnlers.auth.ReturnAuthToken()"

	logger := log.With().Str("fn", op).Logger()
	logger.Debug().Msg("Request for jwt token has been received")

	var userAuthData model.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&userAuthData)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to decode request body")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}
	logger.Debug().Msg("Request body decoded")

	err = validator.New().Struct(userAuthData)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to validate request body")

		w.WriteHeader(http.StatusBadRequest) // 400
		render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
		return
	}

	userID, err := h.postAuth.GetUserID(userAuthData)
	if err != nil {
		if err == db.ErrIncorrectUserPassword {
			logger.Error().Err(err)

			w.WriteHeader(http.StatusUnauthorized) // 401
			render.JSON(w, r, model.ReturnErrResp("Неверный запрос."))
			return
		}
		logger.Error().Err(err).Msg("Failed to get data from the database")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	token, err := token.CreateJWT(userID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create jwt token")

		w.WriteHeader(http.StatusInternalServerError) // 500
		render.JSON(w, r, model.ReturnErrResp("Внутренняя ошибка сервера."))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, model.AuthResponse{Token: token})

}
