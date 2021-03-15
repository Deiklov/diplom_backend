package models

import "time"

type (
	User struct {
		ID        string     `gorm:"primary_key" json:"id" `
		Name      string     `json:"name" valid:"optional,printableascii"`
		Email     string     `json:"email" valid:"required,email"`
		CreatedAt *time.Time `json:"created_at" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
		Password  string     `json:"password" valid:"required,ascii" `
	}
	Profile struct {
		ID         string     `gorm:"primary_key" json:"id"`
		UserID     string     `json:"user_id"`
		Age        string     `json:"age"`
		AvatarPath string     `json:"avatar_path"`
		CreatedAt  *time.Time `json:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at"`
		DeletedAt  *time.Time `json:"deleted_at"`
	}
	AuthData struct {
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"required,ascii"`
	}
	Company struct {
		ID          string `json:"id"`
		Name        string `json:"name" valid:"required,ascii"`
		Year        uint32 `json:"founded_at" valid:"optional,int"`
		Description string `json:"description" valid:"optional,ascii"`
	}
)

func (User) TableName() string { return "users" }
