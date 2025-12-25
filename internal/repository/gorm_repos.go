package repository

import (
	"context"
	"errors"

	"flexi-redirector/internal/models"

	"gorm.io/gorm"
)

type GormRepositories struct {
	ShortURLs ShortURLRepository
}

func NewGormRepositories(gormDB *gorm.DB) GormRepositories {
	return GormRepositories{
		ShortURLs: &gormShortURLRepo{gormDB: gormDB},
	}
}

type gormShortURLRepo struct {
	gormDB *gorm.DB
}

func (repo *gormShortURLRepo) GetBySlug(ctx context.Context, slug string) (models.ShortURL, error) {
	var shortURL models.ShortURL
	getError := repo.gormDB.WithContext(ctx).Where("slug = ?", slug).Take(&shortURL).Error
	if getError != nil {
		if errors.Is(getError, gorm.ErrRecordNotFound) {
			return models.ShortURL{}, ErrNotFound
		}
		return models.ShortURL{}, getError
	}

	return shortURL, nil
}

func (repo *gormShortURLRepo) IncrementViews(ctx context.Context, id uint) error {
	// TODO: check rows affected?
	return repo.gormDB.WithContext(ctx).
		Model(&models.ShortURL{}).
		Where("id = ?", id).
		Update("views", gorm.Expr("views + 1")).Error
}
