package http

import (
	"flexi-redirector/internal/features/countviews"
	"flexi-redirector/internal/repository"
	"log"
	"net/http"

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
	engine.GET("/", func(ginContext *gin.Context) {
		ginContext.Redirect(http.StatusTemporaryRedirect, "/health")
	})

	engine.GET("/favicon.ico", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNoContent)
	})
	engine.GET("/robots.txt", func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, "User-agent: *\nDisallow:\n")
	})
	engine.GET("/sitemap.xml", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNoContent)
	})
	engine.GET("/browserconfig.xml", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNoContent)
	})
	engine.GET("/apple-touch-icon.png", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNoContent)
	})
	engine.GET("/apple-touch-icon-precomposed.png", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNoContent)
	})
	engine.GET("/apple-touch-icon-:size.png", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNoContent)
	})

	engine.GET("/.well-known/*any", func(ginContext *gin.Context) {
		ginContext.Status(http.StatusNotFound)
	})

	// TODO: Add a configurable (on/off) file-path checker. By default, treat a slug containing a '.' as a file request and return 404 to avoid unnecessary DB queries for static assets.

	handlers := NewHandlers(deps)
	engine.GET("/:slug", handlers.Redirect())

	return engine
}
