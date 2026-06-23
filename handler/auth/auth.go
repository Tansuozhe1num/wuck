package auth

import (
	"errors"
	"strings"

	"luangao/store"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore *store.UserStore
}

func NewAuthHandler(userStore *store.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(input RegisterInput) (*store.User, error) {
	if strings.TrimSpace(input.Username) == "" {
		return nil, errors.New("用户名不能为空")
	}
	if strings.TrimSpace(input.Email) == "" {
		return nil, errors.New("邮箱不能为空")
	}
	if len(input.Password) < 6 {
		return nil, errors.New("密码至少需要6个字符")
	}

	email := strings.TrimSpace(input.Email)
	username := strings.TrimSpace(input.Username)

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	return h.userStore.Create(username, email, string(hash))
}

func (h *AuthHandler) Login(input LoginInput) (*store.User, error) {
	if strings.TrimSpace(input.Email) == "" {
		return nil, errors.New("邮箱不能为空")
	}

	email := strings.TrimSpace(input.Email)
	user := h.userStore.FindByEmail(email)
	if user == nil {
		return nil, errors.New("邮箱或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	return user, nil
}

func (h *AuthHandler) GetProfile(userID string) (*store.User, error) {
	user := h.userStore.FindByID(userID)
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (h *AuthHandler) RecordClick(userID string, record store.ClickRecord) error {
	return h.userStore.AddClick(userID, record)
}
