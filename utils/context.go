package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractUserClaims(claims jwt.MapClaims) (userID uint, role string, phone string, ok bool) {
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", "", false
	}
	roleStr, _ := claims["role"].(string)
	phoneStr, _ := claims["phone"].(string)
	return uint(id), roleStr, phoneStr, true
}

func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}
