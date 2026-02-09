package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"flexi-redirector/internal/models"

	"gorm.io/gorm"
)

type GormRepositories struct {
	ShortURLs      ShortURLRepository
	shortURLSchema ShortURLSchema
}

func NewGormRepositories(gormDB *gorm.DB, shortURLSchema ShortURLSchema) GormRepositories {
	return GormRepositories{
		ShortURLs: &gormShortURLRepo{gormDB: gormDB, shortURLSchema: shortURLSchema},
	}
}

type gormShortURLRepo struct {
	gormDB         *gorm.DB
	shortURLSchema ShortURLSchema
}

func (repo *gormShortURLRepo) GetBySlug(ctx context.Context, slug string) (models.ShortURL, error) {
	var shortURL models.ShortURL

	selectClause := repo.shortURLSelectClause()
	getError := repo.gormDB.WithContext(ctx).
		Table(repo.shortURLSchema.TableName).
		Select(selectClause).
		Where(repo.shortURLSchema.ColSlug+" = ?", slug).
		Take(&shortURL).Error
	if getError != nil {
		if errors.Is(getError, gorm.ErrRecordNotFound) {
			return models.ShortURL{}, ErrNotFound
		}
		return models.ShortURL{}, getError
	}

	return shortURL, nil
}

func (repo *gormShortURLRepo) IncrementViews(ctx context.Context, id uint) error {
	viewsColumn, ok := optionalColumn(repo.shortURLSchema.ColViews)
	if !ok {
		return errors.New("views column is disabled for this schema")
	}

	// TODO: check rows affected?
	return repo.gormDB.WithContext(ctx).
		Table(repo.shortURLSchema.TableName).
		Where(repo.shortURLSchema.ColId+" = ?", id).
		Update(viewsColumn, gorm.Expr(viewsColumn+" + 1")).Error
}

func (repo *gormShortURLRepo) shortURLSelectClause() string {
	selectParts := make([]string, 0, 8)

	selectParts = append(selectParts, fmt.Sprintf("%s as id", repo.shortURLSchema.ColId))
	selectParts = append(selectParts, fmt.Sprintf("%s as slug", repo.shortURLSchema.ColSlug))
	selectParts = append(selectParts, fmt.Sprintf("%s as long_url", repo.shortURLSchema.ColLongURL))

	if viewsColumn, ok := optionalColumn(repo.shortURLSchema.ColViews); ok {
		selectParts = append(selectParts, fmt.Sprintf("%s as views", viewsColumn))
	}
	if createdByColumn, ok := optionalColumn(repo.shortURLSchema.ColCreatedBy); ok {
		selectParts = append(selectParts, fmt.Sprintf("%s as created_by", createdByColumn))
	}
	if nameColumn, ok := optionalColumn(repo.shortURLSchema.ColName); ok {
		selectParts = append(selectParts, fmt.Sprintf("%s as name", nameColumn))
	}
	if createdAtColumn, ok := optionalColumn(repo.shortURLSchema.ColCreatedAt); ok {
		selectParts = append(selectParts, fmt.Sprintf("%s as created_at", createdAtColumn))
	}
	if updatedAtColumn, ok := optionalColumn(repo.shortURLSchema.ColUpdatedAt); ok {
		selectParts = append(selectParts, fmt.Sprintf("%s as updated_at", updatedAtColumn))
	}

	return strings.Join(selectParts, ", ")
}

func optionalColumn(columnName string) (string, bool) {
	trimmed := strings.TrimSpace(columnName)
	if trimmed == "" || trimmed == "-" {
		return "", false
	}
	return trimmed, true
}
