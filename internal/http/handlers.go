package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"flexi-redirector/internal/features/countviews"
	"flexi-redirector/internal/repository"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	shortURLs  repository.ShortURLRepository
	countViews *countviews.Feature
	logger     *log.Logger
}

func NewHandlers(deps Deps) *Handlers {
	logger := deps.Logger
	if logger == nil {
		logger = log.Default()
	}
	return &Handlers{shortURLs: deps.ShortURLs, countViews: deps.CountViews, logger: logger}
}

func (handlers *Handlers) Redirect() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		slug := ginContext.Param("slug")
		if slug == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Slug parameter is required"})
			return
		}

		shortURL, getError := handlers.shortURLs.GetBySlug(ginContext.Request.Context(), slug)
		if getError != nil {
			if errors.Is(getError, repository.ErrNotFound) {
				ginContext.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
				return
			}
			handlers.logger.Printf("Database error: %v", getError)
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if handlers.countViews != nil && handlers.countViews.Enabled() {
			shortURLID := shortURL.ID
			if !handlers.countViews.CountAsync() {
				if incrementError := handlers.shortURLs.IncrementViews(ginContext.Request.Context(), shortURLID); incrementError != nil {
					handlers.logger.Printf("Failed to update view count: %v", incrementError)
				}
			} else {
				asyncTimeout := handlers.countViews.AsyncTimeout()
				go func() {
					timeoutContext, cancel := context.WithTimeout(context.Background(), asyncTimeout)
					defer cancel()
					if incrementError := handlers.shortURLs.IncrementViews(timeoutContext, shortURLID); incrementError != nil {
						handlers.logger.Printf("Failed to update view count: %v", incrementError)
					}
				}()
			}
		}

		ginContext.Redirect(http.StatusFound, shortURL.LongURL)
	}
}
