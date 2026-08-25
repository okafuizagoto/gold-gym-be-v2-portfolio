package middleware

import (
	"fmt"
	"strings"

	goldAuthEntity "gold-gym-be/internal/entity/auth/v2"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckUniqueRequest(c *gin.Context) {

	// err := h.service.ValidateHeaders(c.Request.Header)
	// if err != nil {

	// 	var appErr *apperror.AppError
	// 	if errors.As(err, &appErr) {
	// 		c.AbortWithStatusJSON(400, gin.H{
	// 			"responseCode":    appErr.Code,
	// 			"responseMessage": appErr.Message,
	// 		})
	// 		return
	// 	}

	// 	c.AbortWithStatusJSON(500, gin.H{
	// 		"responseCode":    "5000000",
	// 		"responseMessage": "Internal error",
	// 	})
	// 	return
	// }
	fmt.Println("Middleware: CheckUniqueRequest")

	c.Next()
}

func (h *Handler) ValidateToken(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "AUTH_HEADER_MISSING"})
		c.Abort()
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"error": "INVALID_AUTH_FORMAT"})
		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return goldAuthEntity.JwtSecret, nil
	})

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				c.JSON(401, gin.H{"error": "TOKEN_EXPIRED"})
				c.Abort()
				return
			}
		}

		c.JSON(401, gin.H{"error": "INVALID_TOKEN"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		user := int(claims["sub"].(float64))
		c.Set("user", user)

		// role dari claim; token lama tanpa claim role dianggap SELLER
		role := "SELLER"
		if r, ok := claims["role"].(string); ok && r != "" {
			role = r
		}
		c.Set("role", role)

		c.Next()
		return
	}

	c.JSON(401, gin.H{"error": "UNAUTHORIZED"})
	c.Abort()
}
