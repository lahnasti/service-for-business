package repository

import (
	"context"
	"fmt"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
)

func (db *DBstorage) GetReviewsByAuthorAndTender(tenderID int, username string, organizationID int) ([]models.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение userID по имени пользователя (CreatorUsername)
	var userID int
	err := db.conn.WithContext(ctx).
		Table("employee").
		Select("id").
		Where("username = ?", username).
		Scan(&userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Проверка: является ли пользователь ответственным за организацию
	var count int64
	err = db.conn.WithContext(ctx).
		Table("organization_responsible").
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("failed to check user responsibilities: %w", err)
	}

	// Если count == 0, то пользователь не является ответственным за организацию
	if count == 0 {
		return nil, fmt.Errorf("user does not have permission to view reviews")
	}

	// Получение отзывов на предложения, созданные автором, для указанного тендера
	var reviews []models.Review
	err = db.conn.WithContext(ctx).
		Table("reviews").
		Joins("JOIN bid ON reviews.bid_id = bid.id").
		Where("bid.tender_id = ? AND bid.creator_username = ? AND bid.organization_id = ?", tenderID, username, organizationID).
		Find(&reviews).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (db *DBstorage) AddFeedback(reviews models.Review, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Проверка прав пользователя
	hasPermission, err := db.CheckUserPermissionForTender(reviews.BidID, username)
	if err != nil {
		return fmt.Errorf("failed to check user permissions: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf("user does not have permission to add feedback")
	}

	err = db.conn.WithContext(ctx).Table("reviews").Create(&reviews).Error
	if err != nil {
		return fmt.Errorf("failed to add feedback: %w", err)
	}

	return nil
}
