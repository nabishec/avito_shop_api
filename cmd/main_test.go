package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nabishec/avito_shop_api/cmd/db_connection"
	"github.com/nabishec/avito_shop_api/internal/model"
)

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func testGetAuthTokenUser(t *testing.T, s *Server) string {
	reqBody := map[string]string{
		"username": "testuser1",
		"password": "testpass1",
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("failed to marshal request body")
	}

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, response.Code)

	var authResp model.AuthResponse
	err = json.Unmarshal(response.Body.Bytes(), &authResp)
	if err != nil {
		t.Fatal("failed to unmarshal response body")
	}
	return authResp.Token
}

func testAddUser(t *testing.T, s *Server) string {
	reqBody := map[string]string{
		"username": "testuser2",
		"password": "testpass2",
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("failed to marshal request body")
	}

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, response.Code)

	var authResp model.AuthResponse
	err = json.Unmarshal(response.Body.Bytes(), &authResp)
	if err != nil {
		t.Fatal("failed to unmarshal response body")
	}
	return reqBody["username"]
}

func TestBuyItem(t *testing.T) {
	err := LoadEnv()
	if err != nil {
		t.Error("Don't found config")
	}
	dbConnection, err := db_connection.NewDatabaseConnection()
	if err != nil {
		t.Error("Failed init database")
	}
	s := CreateNewServer(dbConnection.DB)
	s.MountHandlers()
	token := testGetAuthTokenUser(t, s)

	t.Run("Correct Buying", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/cup", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusOK, response.Code)
	})

	t.Run("Not authorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/cup", nil)
		req.Header.Set("Authorization", "Bearer "+"sdfgdgfddfgddd")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Incorrect Buying", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/capitan_america", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Incorrect Buying", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/buy/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusBadRequest, response.Code)
	})

}

func TestSendCoin(t *testing.T) {
	err := LoadEnv()
	if err != nil {
		t.Error("Don't found config")
	}
	dbConnection, err := db_connection.NewDatabaseConnection()
	if err != nil {
		t.Error("Failed init database")
	}
	s := CreateNewServer(dbConnection.DB)
	s.MountHandlers()
	token := testGetAuthTokenUser(t, s)
	toUser := testAddUser(t, s)

	t.Run("Correct Sending", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"toUser": toUser,
			"amount": 10,
		}
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal("failed to marshal request body")
		}

		req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonReq))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusOK, response.Code)
	})

	t.Run("Not authorized", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"toUser": toUser,
			"amount": 10,
		}
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal("failed to marshal request body")
		}

		req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonReq))
		req.Header.Set("Authorization", "Bearer "+"jddsakdjaksd")
		req.Header.Set("Content-Type", "application/json")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Incorrect Buying", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"toUser": toUser,
			"amount": -10,
		}
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal("failed to marshal request body")
		}

		req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonReq))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusBadRequest, response.Code)
	})

}

func TestGetInfo(t *testing.T) {
	err := LoadEnv()
	if err != nil {
		t.Error("Don't found config")
	}
	dbConnection, err := db_connection.NewDatabaseConnection()
	if err != nil {
		t.Error("Failed init database")
	}
	s := CreateNewServer(dbConnection.DB)
	s.MountHandlers()
	token := testGetAuthTokenUser(t, s)

	t.Run("Correct Buying", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusOK, response.Code)
	})

	t.Run("Not authorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+"sdfgdgfddfgddd")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusUnauthorized, response.Code)
	})
}

func TestAuth(t *testing.T) {
	err := LoadEnv()
	if err != nil {
		t.Error("Don't found config")
	}
	dbConnection, err := db_connection.NewDatabaseConnection()
	if err != nil {
		t.Error("Failed init database")
	}
	s := CreateNewServer(dbConnection.DB)
	s.MountHandlers()

	t.Run("Correct Sending", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser2",
			"password": "testpass2",
		}
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal("failed to marshal request body")
		}

		req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusOK, response.Code)
	})

	t.Run("Not authorized", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser2",
			"password": 45327645,
		}
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal("failed to marshal request body")
		}

		req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusBadRequest, response.Code)
	})
}
