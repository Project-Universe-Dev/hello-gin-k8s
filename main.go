package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ビルド時に ldflags で値が埋め込まれる変数
var (
	Version   = "dev"
	GitCommit = "unknown"
)

// getEnv はデフォルト値付きで環境変数を取得
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupRouter はルーターを初期化して返す（本番・テスト共通）
func setupRouter() *gin.Engine {
	r := gin.Default()

	// 環境変数の取得
	appEnv := getEnv("APP_ENV", "development")

	// VERSION と GIT_COMMIT は環境変数があればそれを優先、なければビルド時の値
	version := getEnv("VERSION", Version)
	gitCommit := getEnv("GIT_COMMIT", GitCommit)

	// ルートエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from Gin!")
	})

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// バージョン情報エンドポイント
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"commit":  gitCommit,
		})
	})

	// 環境情報エンドポイント
	r.GET("/env", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"environment": appEnv,
		})
	})

	return r
}

func main() {
	// ルーター初期化
	r := setupRouter()

	// サーバー起動（ポート 8080）
	r.Run(":8080")
}
