package model

type AuthRequest struct {
	Name     string `json:"name" db:"name" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ReturnErrResp(errMsg string) ErrorResponse {
	return ErrorResponse{
		Error: errMsg,
	}
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"string" db:"to_user" validate:"required"`
	Amount int    `json:"integer" db:"to_user" validate:"required"`
}

// TODO: add two model for req
