package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	DeleteByRoomIDTx(ctx context.Context, tx *gorm.DB, id int) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (*reviewRepository) DeleteByRoomIDTx(ctx context.Context, tx *gorm.DB, id int) error {
	err := tx.WithContext(ctx).Where("room_id = ?", id).Delete(&models.Review{}).Error
	if err != nil {
		return err
	}
	return nil
}
