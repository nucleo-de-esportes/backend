package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

func AuthUser(supabase *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token ausente ou inválido",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := supabase.Auth.User(c, token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido ou expirado",
			})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)

		if userType, ok := user.UserMetadata["user_type"].(string); ok {
			c.Set("user_type", userType)
		}

		c.Next()

	}

}
