package model

//TODO: CHECK db teg

// swagger:model AuthRequest
type AuthRequest struct {
	Name     string `json:"username" db:"name" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}

// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"errors"`
}

func ReturnErrResp(errMsg string) ErrorResponse {
	return ErrorResponse{
		Error: errMsg,
	}
}

type AuthResponse struct {
	Token string `json:"token"`
}

// swagger:model SendCoinRequest
type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required"`
}

// swagger:model Item
type Item struct {
	Type     string `json:"type" db:"type"`
	Quantity int    `json:"quantity" db:"quantity"`
}

// swagger:model Received
type Received struct {
	FromUser string `json:"fromUser" db:"name"`
	Amount   int    `json:"amount" db:"amount"`
}

// swagger:model Sent
type Sent struct {
	ToUser string `json:"toUser" db:"name"`
	Amount int    `json:"amount" db:"amount"`
}

// swagger:model CoinHistory
type CoinHistory struct {
	Received []Received `json:"received" db:"name"`
	Sent     []Sent     `json:"sent" db:"amount"`
}

// swagger:model InfoResponse
type InfoResponse struct {
	Coins       int         `json:"coins" db:"coins_number"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}
