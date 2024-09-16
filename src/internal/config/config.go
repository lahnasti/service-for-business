package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr             string
	MPath            string
	DebugFlag        bool
	PostgresHost     string
	PostgresPort     string
	PostgresUsername string
	PostgresPass     string
	PostgresDBName   string
}

// Константы по умолчанию
const (
	defaultAddr             = ":8080"
	defaultMigratePath      = "migrations"
	defaultPostgresHost     = "db"
	defaultPostgresPort     = "5432"
	defaultPostgresUsername = "nastya"
	defaultPostgresPass     = "pgspgs"
	defaultPostgresDBName   = "avito"
)

// Функция обработки флагов запуска
func ReadConfig() Config {
	var addr string
	var migratePath string
	debug := flag.Bool("debug", false, "Enable debug logger level")
	flag.StringVar(&addr, "addr", defaultAddr, "Server address")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "Path to migrations")
	flag.Parse()

	if temp := os.Getenv("SERVER_ADDRESS"); temp != "" {
		addr = temp
	}
	if temp := os.Getenv("MIGRATE_PATH"); temp != "" {
		migratePath = temp
	}

	// PostgreSQL environment variables or defaults
	postgresHost := getEnv("POSTGRES_HOST", defaultPostgresHost)
	postgresPort := getEnv("POSTGRES_PORT", defaultPostgresPort)
	postgresUsername := getEnv("POSTGRES_USERNAME", defaultPostgresUsername)
	postgresPass := getEnv("POSTGRES_PASSWORD", defaultPostgresPass)
	postgresDBName := getEnv("POSTGRES_DATABASE", defaultPostgresDBName)

	return Config{
		Addr:             addr,
		MPath:            migratePath,
		DebugFlag:        *debug,
		PostgresHost:     postgresHost,
		PostgresPort:     postgresPort,
		PostgresUsername: postgresUsername,
		PostgresPass:     postgresPass,
		PostgresDBName:   postgresDBName,
	}
}

// Функция для получения переменной окружения или значения по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
