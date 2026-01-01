package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMain はテスト全体の初期化
func TestMain(m *testing.M) {
	// テストモードに設定
	gin.SetMode(gin.TestMode)

	// テスト実行
	code := m.Run()

	os.Exit(code)
}

// TestRootEndpoint はルートパスのテスト
func TestRootEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello from Gin!", w.Body.String())
}

// TestHealthEndpoint はヘルスチェックエンドポイントのテスト
func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

// TestVersionEndpoint はバージョンエンドポイントのテスト
func TestVersionEndpoint(t *testing.T) {
	// テスト用に環境変数を設定
	os.Setenv("VERSION", "1.2.3")
	os.Setenv("GIT_COMMIT", "abc123")
	defer func() {
		os.Unsetenv("VERSION")
		os.Unsetenv("GIT_COMMIT")
	}()

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/version", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3", response["version"])
	assert.Equal(t, "abc123", response["commit"])
}

// TestEnvEndpoint は環境情報エンドポイントのテスト
func TestEnvEndpoint(t *testing.T) {
	os.Setenv("APP_ENV", "testing")
	defer os.Unsetenv("APP_ENV")

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/env", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "testing", response["environment"])
}

// TestGetEnv はgetEnv関数のテスト（ユーティリティ関数のテスト）
func TestGetEnv(t *testing.T) {
	// 環境変数が設定されている場合
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	assert.Equal(t, "test_value", result)

	// 環境変数が設定されていない場合
	result = getEnv("NON_EXISTENT_VAR", "default")
	assert.Equal(t, "default", result)
}
