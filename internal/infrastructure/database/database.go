package database

import (
	"fmt"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.AuditLog{},
		&entity.Store{},
		&entity.Stock{},
		&entity.StockEntry{},
		&entity.StockHistory{},
		&entity.Client{},
		&entity.ClientAddress{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create default admin role if not exists
	var adminRole entity.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == gorm.ErrRecordNotFound {
		adminRole = entity.Role{
			Name: "admin",
			Permissions: entity.GormPermissionSlice{
				entity.UserCreate,
				entity.UserRead,
				entity.UserUpdate,
				entity.UserDelete,
				entity.RoleCreate,
				entity.RoleRead,
				entity.RoleUpdate,
				entity.RoleDelete,
				entity.AuditLogRead,
				entity.ModuleIntegrate,

				// Store permissions
				entity.StoreCreate,
				entity.StoreRead,
				entity.StoreUpdate,
				entity.StoreDelete,

				// Stock permissions
				entity.StockRead,
				entity.StockUpdate,
				entity.StockEntryCreate,
				entity.StockEntryRead,

				// Client permissions
				entity.ClientCreate,
				entity.ClientRead,
				entity.ClientUpdate,
				entity.ClientDelete,
				entity.ClientAddressCreate,
				entity.ClientAddressRead,
				entity.ClientAddressUpdate,
				entity.ClientAddressDelete,
				entity.ClientDebtRead,
				entity.ClientDebtUpdate,
				entity.ClientLoyaltyRead,
				entity.ClientLoyaltyUpdate,
			},
		}

		if err := db.Create(&adminRole).Error; err != nil {
			return nil, fmt.Errorf("failed to create admin role: %w", err)
		}

		// Create default admin user if not exists
		var adminUser entity.User
		if err := db.Where("username = ?", "admin").First(&adminUser).Error; err == gorm.ErrRecordNotFound {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash admin password: %w", err)
			}

			adminUser = entity.User{
				Username: "admin",
				Email:    "admin@example.com",
				Password: string(hashedPassword),
				RoleID:   adminRole.ID,
				Status:   entity.StatusActive,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				return nil, fmt.Errorf("failed to create admin user: %w", err)
			}
		}
	}

	return db, nil
}
