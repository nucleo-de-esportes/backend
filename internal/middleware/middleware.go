package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
)

func AuthUser(c *gin.Context) {

	var tokenString string

	cookieToken, err := c.Cookie("Authorization")
	if err == nil && cookieToken != "" {
		tokenString = cookieToken
	} else {
		//Se nao tem cookie, pega o token do header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token não fornecido",
				"code":  "TOKEN_MISSING",
			})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {

			tokenString = authHeader
		}
	}

	user, err := ValidateTokenAndGetUser(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
			"code":  "AUTHENTICATION_FAILED",
		})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Set("user_id", user.User_id.String())
	c.Set("user_type", user.User_type)

	c.Next()
}

func ValidateTokenAndGetUser(tokenString string) (*model.User, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			return nil, jwt.ErrInvalidKey
		}

		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	var userID string
	if subjectClaim, exists := claims["subject"]; exists {
		if subStr, ok := subjectClaim.(string); ok {
			userID = subStr
		}
	} else if subClaim, exists := claims["sub"]; exists {
		if subStr, ok := subClaim.(string); ok {
			userID = subStr
		}
	}

	if userID == "" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	var user model.User
	result := repository.DB.First(&user, "user_id = ?", userID)

	if result.Error != nil {
		return nil, result.Error
	}

	if user.User_id.IsNil() {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &user, nil
}
