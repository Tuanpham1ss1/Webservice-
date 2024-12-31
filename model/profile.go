package model

import "time"

type Profile struct {
	ID               uint      `json:"id"`
	UserID           uint      `json:"userId"`
	Avatar           *string   `json:"avatar"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Address          string    `json:"address"`
	Code             string    `json:"code"`
	DateOfBirth      string    `json:"dateOfBirth"`
	IsCreatedAccount bool      `json:"isCreatedAccount"`
	OtherInfo        *string   `json:"otherInfo"`
	User             User      `swaggerignore:"true" json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt        time.Time `json:"createdAt" swaggerignore:"true"`
	UpdatedAt        time.Time `json:"updatedAt" swaggerignore:"true"`
	DeletedAt        time.Time `json:"deletedAt" swaggerignore:"true"`
}
