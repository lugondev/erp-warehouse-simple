package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo entity.UserRepository
}

func NewUserUseCase(repo entity.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: repo}
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
	RoleID   uint
}

type UpdateUserInput struct {
	ID       uint
	Username string
	Email    string
	RoleID   uint
	Status   entity.UserStatus
}

type ChangePasswordInput struct {
	UserID      uint
	OldPassword string
	NewPassword string
}

func (uc *UserUseCase) CreateUser(input *CreateUserInput) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   input.RoleID,
		Status:   entity.StatusActive,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) UpdateUser(input *UpdateUserInput) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	user.Username = input.Username
	user.Email = input.Email
	user.RoleID = input.RoleID
	user.Status = input.Status

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Delete(id uint) error {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Role != nil && user.Role.Name == "admin" {
		return errors.New("cannot delete admin user")
	}

	return uc.userRepo.Delete(id)
}

func (uc *UserUseCase) ChangePassword(input *ChangePasswordInput) error {
	user, err := uc.userRepo.FindByID(input.UserID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdatePassword(user.ID, string(hashedPassword))
}

func (uc *UserUseCase) ValidateCredentials(email, password string) (*entity.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive() {
		return nil, errors.New("user account is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := uc.userRepo.UpdateLastLogin(user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) GetUserByID(id uint) (*entity.User, error) {
	return uc.userRepo.FindByID(id)
}

func (uc *UserUseCase) ListUsers() ([]entity.User, error) {
	return uc.userRepo.List()
}

func (uc *UserUseCase) UpdateRefreshToken(userID uint, token string, expiry time.Time) error {
	return uc.userRepo.UpdateRefreshToken(userID, token, expiry)
}

func (uc *UserUseCase) ValidateRefreshToken(token string) (*entity.User, error) {
	return uc.userRepo.FindByRefreshToken(token)
}

func (uc *UserUseCase) RequestPasswordReset(email string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	token := generateResetToken()
	expiry := time.Now().Add(24 * time.Hour)

	if err := uc.userRepo.UpdatePasswordResetToken(user.ID, token, expiry); err != nil {
		return "", err
	}

	return token, nil
}

func (uc *UserUseCase) ResetPassword(token, newPassword string) error {
	user, err := uc.userRepo.FindByPasswordResetToken(token)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := uc.userRepo.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		return err
	}

	// Clear the reset token
	return uc.userRepo.UpdatePasswordResetToken(user.ID, "", time.Time{})
}

func generateResetToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
