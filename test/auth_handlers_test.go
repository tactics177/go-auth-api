package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tactics177/go-auth-api/config"
	"github.com/tactics177/go-auth-api/internal/handlers"
	"github.com/tactics177/go-auth-api/internal/models"
	"github.com/tactics177/go-auth-api/internal/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestMain(m *testing.M) {
	config.ConnectDB()
	os.Exit(m.Run())
}

func clearTestDB() {
	_ = config.DB.Collection("users").Drop(context.TODO())
	_ = config.DB.Collection("refresh_tokens").Drop(context.TODO())
	_ = config.DB.Collection("blacklisted_tokens").Drop(context.TODO())
}

func TestRegister_Success(t *testing.T) {
	clearTestDB()

	router := setupRouter()
	router.POST("/register", handlers.Register)

	user := models.User{
		Name:     "Jane Doe1",
		Email:    "jane2@example.com",
		Password: "securepass123",
	}
	body, _ := json.Marshal(user)

	services.RegisterUserFn = func(user *models.User) error {
		return nil
	}

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestLogin_ValidationError(t *testing.T) {
	router := setupRouter()
	router.POST("/login", handlers.Login)

	body := []byte(`{"email":"", "password":""}`)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
