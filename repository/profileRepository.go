package repository

import (
	"gorm.io/gorm"
	"strings"
	"test1/model"
)

type ProfileRepository interface {
	GetProfile(db *gorm.DB, profileId uint) (*model.Profile, error)
	GetProfileByUsername(db *gorm.DB, username string) (*model.Profile, error)
	CreateProfile(db *gorm.DB, profile *model.Profile) error
	GetProfileByPhoneOrEmail(db *gorm.DB, phone string, email string) (*model.Profile, error)
}
type profileRepository struct {
}

func (p *profileRepository) GetProfile(db *gorm.DB, profileId uint) (*model.Profile, error) {
	var profile model.Profile
	err := db.Preload("User").First(&profile, profileId).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
func (p *profileRepository) GetProfileByUsername(db *gorm.DB, username string) (*model.Profile, error) {
	var profile model.Profile
	err := db.
		Joins("JOIN users ON profiles.user_id = users.id").
		Where("users.username = ?", username).
		Preload("User").
		First(&profile).Error

	if err != nil {
		return nil, err
	}

	return &profile, nil
}
func (p *profileRepository) CreateProfile(db *gorm.DB, profile *model.Profile) error {
	if err := db.Create(profile).Error; err != nil {
		return err
	}
	return nil
}
func (p *profileRepository) GetProfileByPhoneOrEmail(db *gorm.DB, phone string, email string) (*model.Profile, error) {
	var profile model.Profile
	email = strings.ToLower(email)
	temp := db.Where("LOWER(email) = ? OR phone = ?", email, phone)
	if err := temp.First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{}
}
