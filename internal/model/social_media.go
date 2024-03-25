package model

import (
	"time"

	"gorm.io/gorm"
)

type SocialMedia struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url" gorm:"column:url"`
	UserId         uint64    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt
}

type SocialMediaReq struct {
	Name           string `json:"name" validate:"required"`
	SocialMediaUrl string `json:"social_media_url" validate:"required"`
}

type SocialMediaCreateRes struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	UserId         uint64    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type SocialMediaUpdateRes struct {
	ID             uint64    `json:"id"`
	UserId         uint64    `json:"user_id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SocialMediaGetRes struct {
	ID             uint64    	 `json:"id"`
	Name           string    	 `json:"name"`
	SocialMediaUrl string    	 `json:"social_media_url" gorm:"column:url"`
	UserId         uint64    	 `json:"user_id"`
	CreatedAt      time.Time 	 `json:"created_at"`
	UpdatedAt      time.Time 	 `json:"updated_at"`
	User           UserRelation  `json:"User" gorm:"foreignKey:UserId;references:ID"`
}