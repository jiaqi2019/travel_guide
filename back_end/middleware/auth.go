package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"travel_guide/models"
	"travel_guide/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint) (string, error) {
	expiresIn, _ := strconv.ParseInt(getEnv("JWT_EXPIRES_IN", "86400"), 10, 64)
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getEnv("JWT_SECRET_KEY", "")))
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "Authorization header is required"))
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(getEnv("JWT_SECRET_KEY", "")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "Invalid token"))
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// AdminMiddleware checks if the user has admin role
func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "Unauthorized"))
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "User not found"))
			c.Abort()
			return
		}

		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "Admin access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware adds user info to context if available but doesn't require authentication
func OptionalAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(getEnv("JWT_SECRET_KEY", "")), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
