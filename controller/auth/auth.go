package auth

import (
	"net/http"

	"luangao/biu"
	handlerauth "luangao/handler/auth"
	"luangao/store"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authHandler *handlerauth.AuthHandler
}

func NewAuthController(authHandler *handlerauth.AuthHandler) *AuthController {
	return &AuthController{authHandler: authHandler}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var input handlerauth.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		biu.Failed(ctx, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	user, err := c.authHandler.Register(input)
	if err != nil {
		biu.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := biu.GenerateToken(user.ID)
	if err != nil {
		biu.Failed(ctx, http.StatusInternalServerError, "令牌生成失败")
		return
	}

	biu.Success(ctx, gin.H{
		"token": token,
		"user":  sanitizeUser(user),
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input handlerauth.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		biu.Failed(ctx, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	user, err := c.authHandler.Login(input)
	if err != nil {
		biu.Failed(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := biu.GenerateToken(user.ID)
	if err != nil {
		biu.Failed(ctx, http.StatusInternalServerError, "令牌生成失败")
		return
	}

	biu.Success(ctx, gin.H{
		"token": token,
		"user":  sanitizeUser(user),
	})
}

func (c *AuthController) Me(ctx *gin.Context) {
	userID := biu.GetUserID(ctx)
	if userID == "" {
		biu.Failed(ctx, http.StatusUnauthorized, "未认证")
		return
	}

	user, err := c.authHandler.GetProfile(userID)
	if err != nil {
		biu.Failed(ctx, http.StatusNotFound, err.Error())
		return
	}

	biu.Success(ctx, sanitizeUser(user))
}

func (c *AuthController) RecordClick(ctx *gin.Context) {
	userID := biu.GetUserID(ctx)
	if userID == "" {
		biu.Failed(ctx, http.StatusUnauthorized, "未认证")
		return
	}

	var record store.ClickRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		biu.Failed(ctx, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	if err := c.authHandler.RecordClick(userID, record); err != nil {
		biu.Failed(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	biu.Success(ctx, gin.H{"ok": true})
}

func sanitizeUser(u *store.User) gin.H {
	history := u.ClickHistory
	if history == nil {
		history = []store.ClickRecord{}
	}

	return gin.H{
		"id":           u.ID,
		"username":     u.Username,
		"email":        u.Email,
		"clickHistory": history,
		"createdAt":    u.CreatedAt,
	}
}
