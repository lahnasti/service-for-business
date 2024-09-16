package repository

import (
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: Создать тестовую бд для тестирования
type DBstorage struct {
	conn *gorm.DB
}

func NewDB(cfg config.Config) (*DBstorage, error) {
	// Формирование строки подключения к базе данных
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUsername, cfg.PostgresPass, cfg.PostgresDBName, cfg.PostgresPort,
	)

	// Подключение через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DBstorage{
		conn: db,
	}, nil
}
