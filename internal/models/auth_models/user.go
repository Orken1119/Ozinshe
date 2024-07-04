package models

import "context"

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phoneNumber"`
	RoleID      uint   `json:"roleId"`
	CreatedAt   string `json:"createdAt"`
}

type UserRequest struct {
	Email       string   `json:"email"`
	Password    Password `json:"password"`
	PhoneNumber string   `json:"phoneNumber"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Password struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	Email    string   `json:"email"`
	Code     string   `json:"code"`
	Password Password `json:"password"`
}

type UserRepository interface {
	GetUserByEmail(c context.Context, email string) (User, error)
	GetUserByID(c context.Context, userID int) (User, error)
	GetProfile(c context.Context, userID int) (User, error)
	ChangePassword(c context.Context, code string, email string, newPassword string) error
	GetCodeByEmail(c context.Context, email string) (string, error)
	CreatePasswordResetCode(c context.Context, email string, code string) error
	CreateUser(c context.Context, user UserRequest) (int, error)
}
