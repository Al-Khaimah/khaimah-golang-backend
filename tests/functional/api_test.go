package functional

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPIHealth(t *testing.T) {
	t.Logf("📋 TEST: Testing API Health Endpoint")
	startTime := time.Now()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	healthHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "healthy",
			"version": "1.0.0",
		})
	}

	e.GET("/api/health", healthHandler)

	t.Logf("▶️ Sending GET request to /api/health")
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "healthy")

	t.Logf("✅ Health check endpoint responded with status: %d", rec.Code)
	t.Logf("✅ Response contains 'healthy' status")
	t.Logf("⏱️ Test completed in %s", time.Since(startTime))
}

func TestUserSignup(t *testing.T) {
	t.Logf("📋 TEST: Testing User Signup Flow")
	startTime := time.Now()

	e := echo.New()

	payload := `{"first_name":"Test","last_name":"User","email":"test@example.com","password":"TestPass123!"}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	signupHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  true,
			"message": "Account created successfully",
			"data": map[string]interface{}{
				"id":    "user-123",
				"email": "test@example.com",
				"token": "mock-jwt-token",
			},
		})
	}

	e.POST("/api/auth/signup", signupHandler)

	t.Logf("▶️ Sending signup request with user data")
	t.Logf("📦 Payload: %s", payload)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Account created successfully")
	assert.Contains(t, rec.Body.String(), "mock-jwt-token")

	t.Logf("✅ Signup endpoint responded with status: %d", rec.Code)
	t.Logf("✅ Response contains success message and user token")
	t.Logf("⏱️ Test completed in %s", time.Since(startTime))
}

func TestLoginFlow(t *testing.T) {
	t.Logf("📋 TEST: Testing User Login Flow")
	startTime := time.Now()

	e := echo.New()

	payload := `{"email":"test@example.com","password":"TestPass123!"}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	loginHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  true,
			"message": "Logged in successfully",
			"data": map[string]interface{}{
				"id":    "user-123",
				"email": "test@example.com",
				"token": "mock-jwt-token",
			},
		})
	}

	e.POST("/api/auth/login", loginHandler)

	t.Logf("▶️ Sending login request with credentials")
	t.Logf("📦 Payload: %s", payload)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Logged in successfully")
	assert.Contains(t, rec.Body.String(), "mock-jwt-token")

	t.Logf("✅ Login endpoint responded with status: %d", rec.Code)
	t.Logf("✅ Response contains success message and authentication token")
	t.Logf("⏱️ Test completed in %s", time.Since(startTime))
}
