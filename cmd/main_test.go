package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/nabishec/avito_shop_api/cmd/db_connection"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/auth"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/buy"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/info"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/send"
	"github.com/nabishec/avito_shop_api/internal/http_server/middlweare"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/nabishec/avito_shop_api/internal/storage/db"
)

func createTestRouter(t *testing.T) http.Handler {
	testRouter := chi.NewRouter()

	dbConnection, err := db_connection.NewDatabaseConnection()
	if err != nil {
		t.Error("failed connect database")
	}
	t.Log("Database init successful")

	testStorage := db.NewStorage(dbConnection.DB)

	testSendCoin := send.NewSendingCoins(testStorage)
	testGetInformation := info.NewUserInformation(testStorage)
	testBuyItem := buy.NewBuying(testStorage)
	testAuthentication := auth.NewAuth(testStorage)

	testRouter.Group(func(r chi.Router) {
		r.Post("/api/auth", testAuthentication.ReturnAuthToken)
	})

	// Require Authentication
	testRouter.Group(func(r chi.Router) {
		r.Use(middlweare.Auth)
		r.Get("/api/buy/{item}", testBuyItem.BuyingItemByUser)
		r.Post("/api/sendCoin", testSendCoin.SendCoins)
		r.Get("/api/info", testGetInformation.ReturnUserInfo)
	})

	return testRouter
}

func executeRequest(req *http.Request, router http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func testGetAuthTokenUser1(t *testing.T, router http.Handler) string {
	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		t.Error("failed marshal")
	}

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, router)
	checkResponseCode(t, http.StatusOK, response.Code)

	var authResp model.AuthResponse
	json.Unmarshal(response.Body.Bytes(), &authResp)
	return authResp.Token

}

func TestBuyItem(t *testing.T) {
	router := createTestRouter(t)
	token := testGetAuthTokenUser1(t, router)
	t.Run("Correct Buying", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/cup", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response := executeRequest(req, router)
		checkResponseCode(t, http.StatusOK, response.Code)
	})

	t.Run("Not authorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/cup", nil)
		req.Header.Set("Authorization", "Bearer "+"sdfgdgfddfgddd")

		response := executeRequest(req, router)
		checkResponseCode(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Incorrect Buying", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/capitan_america", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response := executeRequest(req, router)
		checkResponseCode(t, http.StatusBadRequest, response.Code)
	})
}
