package dto

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ── Requests ─────────────────────────────────────────────────

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

func (r *CreateUserRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Username) == "" {
		errs["username"] = "username is required"
	}
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "email is required"
	} else if !emailRegex.MatchString(r.Email) {
		errs["email"] = "email format is invalid"
	}
	if strings.TrimSpace(r.FullName) == "" {
		errs["full_name"] = "full_name is required"
	}
	return errs
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	IsActive *bool  `json:"is_active"`
}

func (r *UpdateUserRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Username) == "" {
		errs["username"] = "username is required"
	}
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "email is required"
	} else if !emailRegex.MatchString(r.Email) {
		errs["email"] = "email format is invalid"
	}
	if strings.TrimSpace(r.FullName) == "" {
		errs["full_name"] = "full_name is required"
	}
	return errs
}

// ── Response ─────────────────────────────────────────────────

type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
