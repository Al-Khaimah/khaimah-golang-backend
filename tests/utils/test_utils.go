package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestResponse represents a common API response structure for testing
type TestResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

// PerformRequest performs an HTTP request for testing
func PerformRequest(e *echo.Echo, method, path string, body interface{}) (*httptest.ResponseRecorder, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec, nil
}

// PerformAuthenticatedRequest performs an authenticated HTTP request for testing
func PerformAuthenticatedRequest(e *echo.Echo, method, path string, body interface{}, token string) (*httptest.ResponseRecorder, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec, nil
}

// ParseResponse parses the HTTP response body into a TestResponse
func ParseResponse(rec *httptest.ResponseRecorder) (TestResponse, error) {
	var response TestResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	return response, err
}

// AssertSuccessResponse asserts that the response is successful
func AssertSuccessResponse(t *testing.T, rec *httptest.ResponseRecorder) TestResponse {
	assert.Equal(t, http.StatusOK, rec.Code)

	var response TestResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Status)

	return response
}

// AssertErrorResponse asserts that the response is an error
func AssertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedCode int) TestResponse {
	assert.Equal(t, expectedCode, rec.Code)

	var response TestResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Status)

	return response
}

// CreateRandomEmail generates a random email for testing
func CreateRandomEmail() string {
	return fmt.Sprintf("test.user.%d@example.com", time.Now().UnixNano()&0xFFFFFFFF)
}

// CreateRandomTestData generates random test data
func CreateRandomTestData() map[string]interface{} {
	return map[string]interface{}{
		"first_name": "Test",
		"last_name":  "User",
		"email":      CreateRandomEmail(),
		"password":   "Test123!",
	}
}

// MockResponse creates a mock API response for testing
func MockResponse(status bool, message string, data interface{}) base.Response {
	messageType := "success"
	if !status {
		messageType = "error"
	}

	return base.Response{
		MessageType:  messageType,
		MessageTitle: message,
		Data:         data,
	}
}

// MockErrorResponse creates a mock error response for testing
func MockErrorResponse(message string, errors interface{}) base.Response {
	return base.Response{
		MessageType:  "error",
		MessageTitle: message,
		Errors:       errors,
	}
}

// ExtractToken extracts the token from a login response
func ExtractToken(response TestResponse) (string, error) {
	if !response.Status {
		return "", fmt.Errorf("response has error status")
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data format")
	}

	token, ok := data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

// MockLoginFlow performs the login flow and returns the token
func MockLoginFlow(e *echo.Echo, email, password string) (string, error) {
	loginData := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	rec, err := PerformRequest(e, http.MethodPost, "/api/auth/login", loginData)
	if err != nil {
		return "", err
	}

	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("login failed with status %d", rec.Code)
	}

	response, err := ParseResponse(rec)
	if err != nil {
		return "", err
	}

	return ExtractToken(response)
}

// MockSignupFlow performs the signup flow and returns the token
func MockSignupFlow(e *echo.Echo, userData map[string]interface{}) (string, error) {
	rec, err := PerformRequest(e, http.MethodPost, "/api/auth/signup", userData)
	if err != nil {
		return "", err
	}

	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("signup failed with status %d", rec.Code)
	}

	response, err := ParseResponse(rec)
	if err != nil {
		return "", err
	}

	return ExtractToken(response)
}
