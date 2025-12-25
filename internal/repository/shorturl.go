package repository

import (
	"context"

	"flexi-redirector/internal/models"
)

type ShortURLRepository interface {
	GetBySlug(ctx context.Context, slug string) (models.ShortURL, error)
	IncrementViews(ctx context.Context, id uint) error
}
