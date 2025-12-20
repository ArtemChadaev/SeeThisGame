package repository

import (
	"fmt"
	// Даем локальному пакету псевдоним localMigrate, чтобы не было конфликта
	localMigrate "github.com/ArtemChadaev/SeeThisGame/migrate"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {
	// Теперь используем псевдоним localMigrate для обращения к вашей FS
	d, err := iofs.New(localMigrate.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Здесь migrate — это внешняя библиотека (golang-migrate)
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Выполняем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
