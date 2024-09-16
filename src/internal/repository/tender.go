package repository

import (
	"context"
	"fmt"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
)

func (db *DBstorage) GetAllTenders(serviceType string) ([]models.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tenders []models.Tender

	query := db.conn.WithContext(ctx).Table("tender").Model(&models.Tender{}).Where("status=?", "PUBLISHED")

	if serviceType != "" {
		query = query.Where("service_type=?", serviceType)
	}
	err := query.Order("id ASC").Find(&tenders).Error
	if err != nil {
		return nil, err
	}
	return tenders, nil
}

func (db *DBstorage) GetTendersByUser(username string) ([]models.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tenders []models.Tender

	err := db.conn.WithContext(ctx).Table("tender").
		Where("creator_username = ?", username).
		Order("id ASC").
		Find(&tenders).Error
	if err != nil {
		return nil, err
	}

	return tenders, nil
}

func (db *DBstorage) CreateTender(tender models.Tender) (models.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Устанавливаем начальные значения
	tender.Version = 1
	tender.Status = "CREATED" // Устанавливаем статус

	// Проверка, ответственный ли пользователь за организацию
	var count int64
	err := db.conn.WithContext(ctx).
		Table("organization_responsible").
		Joins("JOIN employee ON employee.id = organization_responsible.user_id").
		Where("organization_responsible.organization_id = ? AND employee.username = ?", tender.OrganizationID, tender.CreatorUsername).
		Count(&count).Error

	if err != nil {
		return models.Tender{}, fmt.Errorf("failed to check responsibility: %w", err)
	}

	if count == 0 {
		return models.Tender{}, fmt.Errorf("user %s is not responsible for organization %d", tender.CreatorUsername, tender.OrganizationID)
	}

	// Создание нового тендера
	tender.Status = "CREATED"
	if err := db.conn.WithContext(ctx).
		Table("tender").
		Create(&tender).Error; err != nil {
		return models.Tender{}, fmt.Errorf("failed to create tender: %w", err)
	}

	return tender, nil
}

func (db *DBstorage) SetTenderStatus(id int, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.conn.WithContext(ctx).
		Table("tender").
		Model(&models.Tender{}).
		Where("id = ?", id).
		Update("status", status)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("no tender found with id %d", id)
	}
	return nil
}

func (db *DBstorage) EditTender(id int, name string, description string) (models.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tender models.Tender
	var currentVersion int

	// Получаем текущую версию тендера
	if err := db.conn.WithContext(ctx).
		Table("tender").
		Where("id = ?", id).
		First(&tender).Error; err != nil {
		return models.Tender{}, err
	}

	currentVersion = tender.Version

	// Сохраняем старую запись в историю
	history := models.TenderHistory{
		TenderID:        tender.ID,
		Name:            tender.Name,
		Description:     tender.Description,
		ServiceType:     tender.ServiceType,
		Status:          tender.Status,
		OrganizationID:  tender.OrganizationID,
		CreatorUsername: tender.CreatorUsername,
		Version:         currentVersion,
	}

	// Сохраняем запись в истории
	err := db.conn.WithContext(ctx).
		Table("tender_history").
		Create(&history).Error
	if err != nil {
		return models.Tender{}, err
	}

	// Обновляем текущую версию
	rowsAffected := db.conn.WithContext(ctx).
		Table("tender").
		Where("id = ? AND version = ?", id, currentVersion).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
			"version":     currentVersion + 1,
		}).RowsAffected

	if rowsAffected == 0 {
		return models.Tender{}, fmt.Errorf("no tender found with id %d or version mismatch", id)
	}

	// Получаем обновленный тендер
	if err := db.conn.WithContext(ctx).
		Table("tender").
		Where("id = ?", id).
		First(&tender).Error; err != nil {
		return models.Tender{}, err
	}

	return tender, nil
}

func (db *DBstorage) RollbackTender(id int, version int) (models.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var tenderH models.TenderHistory
	if err := db.conn.WithContext(ctx).
		Table("tender_history").
		Where("tender_id = ? AND version = ?", id, version).
		First(&tenderH).Error; err != nil {
		return models.Tender{}, fmt.Errorf("version %d for tender %d not found: %v", version, id, err)
	}
	// Обновляем текущий тендер с данными из истории
	err := db.conn.WithContext(ctx).
		Table("tender").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":             tenderH.Name,
			"description":      tenderH.Description,
			"service_type":     tenderH.ServiceType,
			"status":           tenderH.Status,
			"organization_id":  tenderH.OrganizationID,
			"creator_username": tenderH.CreatorUsername,
			"version":          tenderH.Version,
		}).Error
	if err != nil {
		return models.Tender{}, fmt.Errorf("failed to rollback tender: %v", err)
	}
	// Создаем новую запись в истории
	newVersion := tenderH.Version + 1
	newTenderH := models.TenderHistory{
		TenderID:        id,
		Version:         newVersion,
		Name:            tenderH.Name,
		Description:     tenderH.Description,
		ServiceType:     tenderH.ServiceType,
		Status:          tenderH.Status,
		OrganizationID:  tenderH.OrganizationID,
		CreatorUsername: tenderH.CreatorUsername,
	}
	if err := db.conn.WithContext(ctx).
		Table("tender_history").
		Create(&newTenderH).Error; err != nil {
		return models.Tender{}, fmt.Errorf("failed to create new history entry: %w", err)
	}

	var updateTender models.Tender
	if err := db.conn.WithContext(ctx).
		Table("tender").
		Where("id = ?", id).
		First(&updateTender).Error; err != nil {
		return models.Tender{}, fmt.Errorf("failed to fetch updated tender: %v", err)
	}
	return updateTender, nil
}
