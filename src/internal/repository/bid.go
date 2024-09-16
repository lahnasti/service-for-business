package repository

import (
	"context"
	"fmt"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
)

func (db *DBstorage) GetBidsByUser(username string) ([]models.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var bids []models.Bid
	err := db.conn.WithContext(ctx).
		Table("bid").
		Where("creator_username = ?", username).
		Order("id ASC").
		Find(&bids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get bids by user: %w", err)
	}
	return bids, nil
}

func (db *DBstorage) GetBidsForTender(tenderID int) ([]models.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var bids []models.Bid
	err := db.conn.WithContext(ctx).
		Table("bid").
		Where("tender_id = ?", tenderID).
		Find(&bids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get bids for tender: %w", err)
	}
	return bids, nil
}

func (db *DBstorage) CreateBid(bid models.Bid, creatorUsername string) (models.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bid.Version = 1
	bid.Status = "CREATED"

	// Получение userID
	userID, err := db.GetUserIDByUsername(creatorUsername)
	if err != nil {
		return models.Bid{}, fmt.Errorf("failed to get user ID: %w", err)
	}
	//проверка прав пользователя
	if bid.OrganizationID != nil {
		var count int64
		err = db.conn.WithContext(ctx).
			Table("organization_responsible").
			Where("user_id = ? OR organization_id = ?", userID, bid.OrganizationID).
			Count(&count).Error
		if err != nil {
			return models.Bid{}, fmt.Errorf("failed to check responsibility: %w", err)
		}
		if count == 0 {
			return models.Bid{}, fmt.Errorf("user %s does not have permission to create bid", bid.CreatorUsername)
		}
	}
	// Проверка существования тендера и его статуса
	var tenderStatus string
	err = db.conn.WithContext(ctx).
		Table("tender").
		Select("status").
		Where("id = ?", bid.TenderID).
		Scan(&tenderStatus).Error
	if err != nil {
		return models.Bid{}, fmt.Errorf("failed to check tender existence: %w", err)
	}
	if tenderStatus == "" {
		return models.Bid{}, fmt.Errorf("tender not found")
	}
	if tenderStatus != "PUBLISHED" {
		return models.Bid{}, fmt.Errorf("cannot create bid, tender is not in PUBLISHED status")
	}
	//создание нового предложения
	bid.Status = "CREATED"
	if err := db.conn.WithContext(ctx).
		Table("bid").
		Create(&bid).Error; err != nil {
		return models.Bid{}, fmt.Errorf("error creating bid: %w", err)
	}
	return bid, nil
}

func (db *DBstorage) SetBidStatus(id int, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.conn.WithContext(ctx).
		Table("bid").
		Model(&models.Bid{}).
		Where("id = ?", id).
		Update("status", status)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("no bid found with id %d", id)
	}
	return nil
}

func (db *DBstorage) EditBid(id int, name string, description string) (models.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var bid models.Bid
	var currentVersion int
	err := db.conn.WithContext(ctx).
		Table("bid").
		Where("id =?", id).
		First(&bid).Error
	if err != nil {
		return models.Bid{}, fmt.Errorf("failed to get bid: %w", err)
	}
	currentVersion = bid.Version

	// Сохраняем старую запись в историю
	history := models.BidHistory{
		BidID:           bid.ID,
		Name:            bid.Name,
		Description:     bid.Description,
		Status:          bid.Status,
		OrganizationID:  bid.OrganizationID,
		TenderID:        bid.TenderID,
		CreatorUsername: bid.CreatorUsername,
		Version:         currentVersion,
	}
	// Сохраняем запись в истории
	err = db.conn.WithContext(ctx).
		Table("bid_history").
		Create(&history).Error
	if err != nil {
		return models.Bid{}, err
	}
	// Обновляем текущую версию
	rowsAffected := db.conn.WithContext(ctx).
		Table("bid").
		Where("id = ? AND version = ?", id, currentVersion).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
			"version":     currentVersion + 1,
		}).RowsAffected

	if rowsAffected == 0 {
		return models.Bid{}, fmt.Errorf("no bid found with id %d or version mismatch", id)
	}

	if err := db.conn.WithContext(ctx).
		Table("bid").
		Where("id = ?", id).
		First(&bid).Error; err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (db *DBstorage) RollbackBid(id int, version int) (models.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var bidH models.BidHistory
	if err := db.conn.WithContext(ctx).
		Table("bid_history").
		Where("bid_id =? AND version =?", id, version).
		First(&bidH).Error; err != nil {
		return models.Bid{}, fmt.Errorf("version %d for bid %d not found: %v", version, id, err)
	}
	// Восстанавливаем предыдущую версию
	err := db.conn.WithContext(ctx).
		Table("bid").
		Where("id =?", id).
		Updates(map[string]interface{}{
			"name":             bidH.Name,
			"description":      bidH.Description,
			"status":           bidH.Status,
			"organization_id":  bidH.OrganizationID,
			"tender_id":        bidH.TenderID,
			"creator_username": bidH.CreatorUsername,
			"version":          bidH.Version,
		}).Error
	if err != nil {
		return models.Bid{}, fmt.Errorf("error rollback bid: %w", err)
	}
	err = db.conn.WithContext(ctx).
		Table("bid_history").
		Where("bid_id = ? AND version > ?", id, version).
		Delete(&models.BidHistory{}).Error
	if err != nil {
		return models.Bid{}, fmt.Errorf("error deleting history: %w", err)
	}
	var updateBid models.Bid
	if err := db.conn.WithContext(ctx).
		Table("bid").
		Where("id =?", id).
		First(&updateBid).Error; err != nil {
		return models.Bid{}, fmt.Errorf("error getting updated bid: %w", err)
	}
	return updateBid, nil
}

func (db *DBstorage) SubmitDecision(bid int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//проверка прав пользователя
	yes, err := db.CheckUserPermissionForBid(bid, username)
	if err != nil {
		return fmt.Errorf("failed to check user permission: %w", err)
	}
	if !yes {
		return fmt.Errorf("user does not have permission to submit decision")
	}

	//проверка статуса
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

	// Проверка существующих решений с статусом "DECLINED"
	var declinedCount int64
	err = db.conn.WithContext(ctx).
		Table("bid_decision").
		Where("bid_id = ? AND decision_status = ?", bid, "DECLINED").
		Count(&declinedCount).Error
	if err != nil {
		return fmt.Errorf("failed to check for declined decisions: %w", err)
	}
	// Если есть хотя бы одно решение "DECLINED", предложение отклоняется
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

	if submittedCount >= quorum {
		// Логика для согласования предложения
		err = db.conn.WithContext(ctx).
			Table("bid").
			Where("id = ?", bid).
			Update("status", "SUBMITTED").Error
		if err != nil {
			return fmt.Errorf("error updating bid status: %w", err)
		}
		// Закрытие тендера после согласования предложения
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
	//проверка прав пользователя
	yes, err := db.CheckUserPermissionForBid(bid, username)
	if err != nil {
		return fmt.Errorf("failed to check user permission: %w", err)
	}
	if !yes {
		return fmt.Errorf("user does not have permission to submit decision")
	}
	//проверка статуса
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

	// Тендер не закрывается при отклонении предложения, согласно новой логике
	return nil
}

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
