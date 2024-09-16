package repository

import (
	"context"
	"fmt"
	"time"
)

func (db *DBstorage) SubmitDecision(bid int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверка прав пользователя
	yes, err := db.CheckUserPermissionForBid(bid, username)
	if err != nil {
		return fmt.Errorf("failed to check user permission: %w", err)
	}
	if !yes {
		return fmt.Errorf("user does not have permission to submit decision")
	}

	// Проверка статуса
	var currentStatus string
	err = db.conn.WithContext(ctx).
		Table("bid").
		Select("status").
		Where("id = ?", bid).
		Scan(&currentStatus).Error
	if err != nil {
		return fmt.Errorf("failed to get current bid status: %w", err)
	}
	if currentStatus != "PUBLISHED" {
		return fmt.Errorf("bid must be in PUBLISHED status to submit decision")
	}

	// Проверка существующих решений со статусом "DECLINED"
	var declinedCount int64
	err = db.conn.WithContext(ctx).
		Table("bid_decisions").
		Where("bid_id = ? AND decision_status = ?", bid, "DECLINED").
		Count(&declinedCount).Error
	if err != nil {
		return fmt.Errorf("failed to check for declined decisions: %w", err)
	}
	if declinedCount > 0 {
		err = db.conn.WithContext(ctx).
			Table("bid").
			Where("id = ?", bid).
			Update("status", "DECLINED").Error
		if err != nil {
			return fmt.Errorf("failed to update bid status to DECLINED: %w", err)
		}
		return nil
	}

	// Сохраняем новое решение "SUBMITTED"
	err = db.conn.WithContext(ctx).
		Table("bid_decisions").
		Create(map[string]interface{}{
			"bid_id":          bid,
			"username":        username,
			"decision_status": "SUBMITTED",
		}).Error
	if err != nil {
		return fmt.Errorf("failed to save submitted decision: %w", err)
	}

	// Получаем количество ответственных за организацию
	var responsibleCount int64
	err = db.conn.WithContext(ctx).
		Table("organization_responsible").
		Where("organization_id = (SELECT organization_id FROM bid WHERE id = ?)", bid).
		Count(&responsibleCount).Error
	if err != nil {
		return fmt.Errorf("failed to get responsible count for organization: %w", err)
	}

	// Вычисляем кворум
	quorum := int64(3)
	if responsibleCount < quorum {
		quorum = responsibleCount
	}

	// Получаем количество решений "SUBMITTED"
	var submittedCount int64
	err = db.conn.WithContext(ctx).
		Table("bid_decisions").
		Where("bid_id = ? AND decision_status = ?", bid, "SUBMITTED").
		Count(&submittedCount).Error
	if err != nil {
		return fmt.Errorf("failed to get submitted decisions: %w", err)
	}

	// Проверяем, достигнут ли кворум
	if submittedCount >= quorum {
		// Обновляем статус предложения на "SUBMITTED"
		err = db.conn.WithContext(ctx).
			Table("bid").
			Where("id = ?", bid).
			Update("status", "SUBMITTED").Error
		if err != nil {
			return fmt.Errorf("error updating bid status: %w", err)
		}

		// Закрываем связанный тендер
		err = db.conn.WithContext(ctx).
			Table("tender").
			Where("id = (SELECT tender_id FROM bid WHERE id = ?)", bid).
			Update("status", "CLOSED").Error
		if err != nil {
			return fmt.Errorf("error updating tender status: %w", err)
		}
	}

	return nil
}
func (db *DBstorage) DeclineDecision(bid int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получаем ID тендера, связанного с предложением
	var tenderID int
	err := db.conn.WithContext(ctx).
		Table("bid").
		Select("tender_id").
		Where("id = ?", bid).
		Scan(&tenderID).Error
	if err != nil {
		return fmt.Errorf("failed to get tender ID for bid: %w", err)
	}

	// Проверка прав пользователя
	hasPermission, err := db.CheckUserPermissionForBid(bid, username)
	if err != nil {
		return fmt.Errorf("failed to check user permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf("user does not have permission to submit decision")
	}

	// Проверка статуса
	var currentStatus string
	err = db.conn.WithContext(ctx).
		Table("bid").
		Select("status").
		Where("id = ?", bid).
		Scan(&currentStatus).Error
	if err != nil {
		return fmt.Errorf("failed to get current bid status: %w", err)
	}
	if currentStatus != "PUBLISHED" {
		return fmt.Errorf("bid must be in PUBLISHED status to submit decision")
	}

	// Проверка существующих решений "DECLINED"
	var declinedCount int64
	err = db.conn.WithContext(ctx).
		Table("bid_decisions").
		Where("bid_id = ? AND decision_status = ?", bid, "DECLINED").
		Count(&declinedCount).Error
	if err != nil {
		return fmt.Errorf("failed to check for declined decisions: %w", err)
	}

	// Если есть хотя бы одно решение "DECLINED", сразу отклоняем предложение
	if declinedCount > 0 {
		err = db.conn.WithContext(ctx).
			Table("bid").
			Where("id = ?", bid).
			Update("status", "DECLINED").Error
		if err != nil {
			return fmt.Errorf("failed to update bid status to DECLINED: %w", err)
		}
		return nil
	}

	// Добавляем решение "DECLINED" в таблицу решений
	err = db.conn.WithContext(ctx).
		Table("bid_decisions").
		Create(map[string]interface{}{
			"bid_id":          bid,
			"username":        username,
			"decision_status": "DECLINED",
		}).Error
	if err != nil {
		return fmt.Errorf("failed to save declined decision: %w", err)
	}

	// Отклоняем предложение после первого решения "DECLINED"
	err = db.conn.WithContext(ctx).
		Table("bid").
		Where("id = ?", bid).
		Update("status", "DECLINED").Error
	if err != nil {
		return fmt.Errorf("failed to update bid status to DECLINED: %w", err)
	}

	return nil
}
