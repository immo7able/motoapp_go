package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"motorcycleApp/domain/model"
	"motorcycleApp/utils"
	"net/http"
	"strings"
)

func JWTAuthSecuredMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookie
			}
		}

		if tokenStr == "" {
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, role, phone, ok := utils.ExtractUserClaims(claims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("phone", phone)
		c.Next()
	}
}

func JWTAuthMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.Set("role", nil)
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, roleStr, phone, ok := utils.ExtractUserClaims(claims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		role := model.Role(roleStr)

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("phone", phone)
		c.Next()
	}
}

func AdminOnly(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookie
			}
		}

		if tokenStr == "" {
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, roleStr, phone, ok := utils.ExtractUserClaims(claims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("role", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		role := model.Role(roleStr)

		if role != model.RoleAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("phone", phone)
		c.Next()
	}
}
