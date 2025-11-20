package middleware

import (
	"net/http"
	"strings"

	"esst_sendEmail/internal/pkg/auth"
	"esst_sendEmail/internal/pkg/code"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware JWT 驗證中間件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Authorization header 獲取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 嘗試從 cookie 獲取
			tokenCookie, err := c.Cookie("token")
			if err != nil || tokenCookie == "" {
				c.JSON(http.StatusUnauthorized, code.GetCodeMessage(code.JWTRejected, "Authorization token required"))
				c.Abort()
				return
			}
			authHeader = "Bearer " + tokenCookie
		}

		// 檢查 token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, code.GetCodeMessage(code.JWTRejected, "Invalid token format"))
			c.Abort()
			return
		}

		// 驗證 token
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, code.GetCodeMessage(code.JWTRejected, "Invalid or expired token"))
			c.Abort()
			return
		}

		// 將用戶資訊存儲在 context 中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理員權限驗證中間件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, code.GetCodeMessage(code.PermissionDenied, "Admin permission required"))
			c.Abort()
			return
		}
		c.Next()
	}
}
