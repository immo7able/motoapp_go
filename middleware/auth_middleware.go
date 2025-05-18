package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"motorcycleApp/utils"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(secretKey []byte) gin.HandlerFunc {
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
			c.Set("isAuthenticated", false)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("isAuthenticated", false)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("isAuthenticated", false)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, role, phone, ok := utils.ExtractUserClaims(claims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("isAuthenticated", false)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("phone", phone)
		c.Set("isAuthenticated", true)
		c.Next()
	}
}

func JWTAuthRedirectMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		c.Set("isAuthenticated", false)
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}
}
