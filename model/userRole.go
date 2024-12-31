package model

import (
	"gorm.io/gorm"
	"time"
)

type UserRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProfileId uint           `json:"profileId"`
	Profile   Profile        `json:"profile" gorm:"foreignKey:ProfileId"`
	CreatedAt time.Time      `json:"createdAt" gorm:"default:current_timestamp;"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	UpdatedAt *time.Time     `json:"updatedAt"  gorm:"autoUpdateTime:milli"`
}
