package model

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url" gorm:"column:url"` 
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type PhotoCreateReq struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption" validate:"required"` 
	PhotoUrl string `json:"photo_url" validate:"required"`
}

type PhotoCreateRes struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

}

type PhotoGetRes struct {
	ID        uint64    	`json:"id"`
	Title     string    	`json:"title"`
	Caption   string    	`json:"caption"`
	PhotoUrl  string    	`json:"photo_url" gorm:"column:url"`
	UserId    uint64    	`json:"user_id"`
	CreatedAt time.Time 	`json:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at"`
	User      UserRelation  `json:"User" gorm:"foreignKey:UserId;references:ID"`
}

type PhotoUpdateReq struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption" validate:"required"`
	PhotoUrl string `json:"photo_url" validate:"required"`
}

type PhotoUpdateRes struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserId    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PhotoRelation struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" gorm:"column:url"`
	UserId   uint64 `json:"user_id"`
}