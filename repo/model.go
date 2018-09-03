package repo

import (
	"time"
)

// User is user_auth's model
type User struct {
	ID        string    `json:"id" db:"id"`
	Name	  string    `json:"name" db:"name"`
	Tingkat   string    `json:"tingkat" db:"tingkat"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
}

// UserRole is user_role's model
type UserRole struct {
	RoleID    string    `json:"role_id" db:"id"`
	UserID    string    `json:"id" db:"user_id"`
	Role      int       `json:"role" db:"role"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
