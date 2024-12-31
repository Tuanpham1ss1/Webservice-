package repository

import (
	"gorm.io/gorm"
	"test1/model"
)

type AuthRepository interface {
	LoginGoogle() error
	IsDuplicateEmail(email string) (db *gorm.DB)
}
type authRepository struct {
	db *gorm.DB
}

func (a *authRepository) LoginGoogle() error {
	return nil
}

func (a *authRepository) IsDuplicateEmail(email string) (db *gorm.DB) {
	var user model.User
	return db.Where("email = ?", email).First(&user)
}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}
