package biu

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userIDKey = "userID"

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateToken(userID string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not configured")
	}

	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return userID, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			Failed(ctx, http.StatusUnauthorized, "未提供认证令牌")
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			Failed(ctx, http.StatusUnauthorized, "认证令牌格式错误")
			ctx.Abort()
			return
		}

		userID, err := ParseToken(parts[1])
		if err != nil {
			Failed(ctx, http.StatusUnauthorized, "认证令牌无效或已过期")
			ctx.Abort()
			return
		}

		ctx.Set(userIDKey, userID)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) string {
	id, _ := ctx.Get(userIDKey)
	if id == nil {
		return ""
	}
	return id.(string)
}
