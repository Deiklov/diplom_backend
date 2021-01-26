package models

import "time"

type User struct {
	ID        string     `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Profile struct {
	ID         string     `gorm:"primary_key" json:"id"`
	UserID     int64      `json:"user_id"`
	Age        string     `json:"age"`
	AvatarPath string     `json:"avatar_path"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (User) TableName() string { return "users" }
