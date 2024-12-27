package models

import (
	"time"
)

// User represents a user of the system
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"type:varchar(20);default:'user'" json:"role"` // 'user' or 'admin'
	CreatedAt    time.Time `json:"created_at"`
}

// UserLoginInput represents the input for the login request
type UserLoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
