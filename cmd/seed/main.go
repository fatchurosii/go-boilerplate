package main

import (
	"errors"
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-boilerplate/internal/config"
	"go-boilerplate/internal/database"
	"go-boilerplate/internal/role"
	"go-boilerplate/internal/user"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	if cfg.DatabaseURL == "" {
		logger.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			logger.Error("database close failed", slog.Any("error", err))
		}
	}()

	adminRole, err := seedAdminRole(db)
	if err != nil {
		logger.Error("seed admin role failed", slog.Any("error", err))
		os.Exit(1)
	}

	if err := seedAdminUser(db, adminRole); err != nil {
		logger.Error("seed admin user failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("seed completed")
}

func seedAdminRole(db *gorm.DB) (*role.Role, error) {
	var admin role.Role
	err := db.Where("slug = ?", "admin").First(&admin).Error
	if err == nil {
		return &admin, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	admin = role.Role{
		Name:     "Admin",
		Slug:     "admin",
		IsActive: true,
	}

	return &admin, db.Create(&admin).Error
}

func seedAdminUser(db *gorm.DB, adminRole *role.Role) error {
	var admin user.User
	err := db.Where("username = ?", "admin").First(&admin).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin = user.User{
		RoleID:   adminRole.ID,
		Username: "admin",
		Password: string(password),
		IsActive: true,
	}

	return db.Create(&admin).Error
}
