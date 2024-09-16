package repository

import (
	"context"
	"fmt"
	"time"
)

// вспомогательные функции
func (db *DBstorage) GetUserIDByUsername(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userID int
	err := db.conn.WithContext(ctx).
		Table("employee").
		Where("username = ?", username).
		Pluck("id", &userID).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID by username: %w", err)
	}
	return userID, nil
}

func (db *DBstorage) CheckUserPermissionForTender(tenderID int, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение userID по имени пользователя (CreatorUsername)
	userID, err := db.GetUserIDByUsername(username)
	if err != nil {
		return false, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Проверка: является ли пользователь автором предложения или ответственным за организацию
	var count int64
	err = db.conn.WithContext(ctx).
		Table("bid").
		Joins("JOIN organization_responsible ON bid.organization_id = organization_responsible.organization_id").
		Where("(bid.tender_id = ? AND (bid.creator_username = ? OR organization_responsible.user_id = ?))", tenderID, username, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user permissions: %w", err)
	}

	// Если count > 0, то пользователь имеет права
	return count > 0, nil
}

func (db *DBstorage) CheckUserPermissionForBid(bidID int, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение userID по имени пользователя (CreatorUsername)
	userID, err := db.GetUserIDByUsername(username)
	if err != nil {
		return false, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Проверка: является ли пользователь автором предложения или ответственным за организацию
	var count int64
	err = db.conn.WithContext(ctx).
		Table("bid").
		Joins("JOIN organization_responsible ON bid.organization_id = organization_responsible.organization_id").
		Where("(bid.id = ? AND (bid.creator_username = ? OR organization_responsible.user_id = ?))", bidID, username, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user permissions: %w", err)
	}

	// Если count > 0, то пользователь имеет права
	return count > 0, nil
}

func (db *DBstorage) GetTenderIDByBidID(bidID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tenderID int
	err := db.conn.WithContext(ctx).
		Table("bid").
		Select("tender_id").
		Where("id = ?", bidID).
		Scan(&tenderID).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get tender ID by bid ID: %w", err)
	}

	return tenderID, nil
}
