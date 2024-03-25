package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint64    `json:"id"`
	UserId    uint64    `json:"user_id"`
	PhotoId   uint64    `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type CommentCreateReq struct {
	Message string `json:"message" validate:"required"`
	PhotoId uint64 `json:"photo_id" validate:"required"`
}

type CommentCreateRes struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoId   uint64    `json:"photo_id"`
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentGetRes struct {
	ID        uint64    	`json:"id"`
	Message   string    	`json:"message"`
	UserId    uint64    	`json:"user_id"`
	PhotoId   uint64    	`json:"photo_id"`
	CreatedAt time.Time 	`json:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at"`
	User      UserRelation  `json:"User" gorm:"foreignKey:UserId;references:ID"`
	Photo     PhotoRelation `json:"Photo" gorm:"foreignKey:PhotoId;references:ID"`
}

type CommentUpdateReq struct {
	Message string `json:"message" validate:"required"`
}

type CommentUpdateRes struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoId   uint64    `json:"photo_id"`
	UserId    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}