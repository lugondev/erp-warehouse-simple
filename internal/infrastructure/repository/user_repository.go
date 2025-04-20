package repository

import (
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *UserRepository) List() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateRefreshToken(userID uint, token string, expiry time.Time) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"refresh_token":        token,
		"refresh_token_expiry": expiry,
	}).Error
}

func (r *UserRepository) FindByRefreshToken(token string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").Where("refresh_token = ? AND refresh_token_expiry > ?", token, time.Now()).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdatePasswordResetToken(userID uint, token string, expiry time.Time) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_reset_token": token,
		"reset_token_expiry":   expiry,
	}).Error
}

func (r *UserRepository) FindByPasswordResetToken(token string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").Where("password_reset_token = ? AND reset_token_expiry > ?", token, time.Now()).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *UserRepository) UpdateLastLogin(userID uint) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("last_login", time.Now()).Error
}
