package http

import (
	"log"
	"net/http"

	"flexi-redirector/internal/features/countviews"
	"flexi-redirector/internal/repository"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	ShortURLs  repository.ShortURLRepository
	CountViews *countviews.Feature
	Logger     *log.Logger
}

func NewRouter(deps Deps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // TODO: make configurable
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	engine.GET("/health", func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	handlers := NewHandlers(deps)
	engine.GET("/:slug", handlers.Redirect())

	return engine
}
