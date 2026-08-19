package database

import (
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(db *sqlx.DB) error {
	// 1. Читаем все файлы из папки migrations
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var upMigrations []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upMigrations = append(upMigrations, entry.Name())
		}
	}
	sort.Strings(upMigrations) // Сортируем по имени (000001, 000002...)

	// 2. Получаем текущую версию БД
	currentVersion, err := getCurrentMigrationVersion(db)
	if err != nil {
		return err
	}

	// 3. Применяем новые миграции в транзакции
	for _, filename := range upMigrations {
		// Извлекаем номер версии из имени файла (например, "000002" из "000002_add_users.up.sql")
		parts := strings.Split(filename, "_")
		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid migration file name %s: %w", filename, err)
		}

		// Если эта миграция уже была применена — пропускаем
		if version <= currentVersion {
			continue
		}

		// Читаем содержимое файла миграции
		content, err := migrationFiles.ReadFile("migrations/" + filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Выполняем миграцию в транзакции
		tx, err := db.Beginx()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Обновляем текущую версию в БД
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update schema version for %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func getCurrentMigrationVersion(db *sqlx.DB) (int, error) {
	// Сначала проверяем, создана ли таблица миграций
	var tableExists int
	err := db.Get(&tableExists, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='schema_migrations'")
	if err != nil {
		return 0, err
	}
	if tableExists == 0 {
		return 0, nil
	}

	// Получаем максимальную примененную версию
	var version int
	err = db.Get(&version, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations")
	if err != nil {
		return 0, err
	}
	return version, nil
}
