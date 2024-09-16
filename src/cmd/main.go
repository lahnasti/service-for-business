package main

import (
	"context"
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/config"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/logger"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/server"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/server/routes"
)

func main() {
	fmt.Println("Server starting")

	// Чтение конфигурации
	cfg := config.ReadConfig()
	fmt.Println(cfg)

	// Настройка логгера
	zlog := logger.SetupLogger(cfg.DebugFlag)
	zlog.Debug().Any("config", cfg).Msg("Check cfg value")
	zlog.Debug().Str("migration_path", cfg.MPath).Msg("Path to migrations")

	// Формирование строки подключения к базе данных
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUsername, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDBName,
	)

	// Выполнение миграций
	err := repository.Migrations(dsn, cfg.MPath, zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Init migrations failed")
	}

	// Создание хранилища данных
	dbStorage, err := repository.NewDB(cfg)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Unable to create database storage")
	}

	// Создание сервера
	server := server.New(context.Background(), dbStorage, zlog)

	// Настройка маршрутов
	r := routes.SetupRoutes(server)

	// Запуск сервера
	zlog.Info().Msgf("Starting server on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		zlog.Fatal().Err(err).Msg("Failed to start server")
	}
}

/*func initDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}*/
