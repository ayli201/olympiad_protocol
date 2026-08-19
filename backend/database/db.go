package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"protocol/backend/database/seeders"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func InitDB(appName string) (*sqlx.DB, error) {
	userDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config dir: %w", err)
	}

	// Создаем подпапку для нашего Wails приложения
	appDataDir := filepath.Join(userDir, appName)
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app data dir: %w", err)
	}

	dbPath := filepath.Join(appDataDir, "app.db")

	_, statErr := os.Stat(dbPath)
	isNewDB := os.IsNotExist(statErr)

	// Подключаемся с прагмами оптимизации SQLite (WAL режим, увеличенный таймаут блокировок)
	dsn := fmt.Sprintf("%s?_time_format=sqlite&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)", dbPath)
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}

	// Запускаем миграции сразу после успешного подключения
	if isNewDB {
		if err := RunMigrations(db); err != nil {
			db.Close()
			log.Fatal(err)
			return nil, fmt.Errorf("migration error: %w", err)
		}

		seeders.SeedSchools(db)
		seeders.SeedSettings(db)
		seeders.SeedWinnersQuota(db)
	}

	return db, nil
}
