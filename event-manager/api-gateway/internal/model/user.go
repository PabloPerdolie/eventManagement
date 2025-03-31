package model

import (
	"time"

	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
)

//CREATE TABLE users
//(
//user_id    SERIAL PRIMARY KEY,
//username   VARCHAR(30) NOT NULL,
//password_hash TEXT NOT NULL,
//email VARCHAR NOT NULL,
//is_active BOOLEAN NOT NULL DEFAULT TRUE,
//is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
//role VARCHAR NOT NULL,
//created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
//);

type User struct {
	UserId       int             `json:"id" db:"id"`
	Username     string          `json:"username" db:"username"`
	Email        string          `json:"email" db:"email"`
	PasswordHash string          `json:"-" db:"password_hash"`
	Role         domain.UserRole `json:"role" db:"role"`
	IsActive     bool            `json:"is_active" db:"is_active"`
	IsDeleted    bool            `json:"is_deleted" db:"is_deleted"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}
