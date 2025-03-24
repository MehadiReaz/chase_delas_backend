package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8"`
	Role      string    `json:"role" validate:"oneof=superadmin admin vendor user"`
	CreatedAt time.Time // Automatically set when created
	UpdatedAt time.Time // Automatically updated
}

type UpdateCredentialsRequest struct {
	Email       string `json:"email,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ErrorResponseModel struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccessResponseModel struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
