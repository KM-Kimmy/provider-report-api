package middleware

import (
	utils "provider-report-api/pkg/utility"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is jwt middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		var data interface{}
		var err error
		code = 200
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			code = 400
			err = errors.New("missing token")
			fmt.Println("missing token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			// ctx.Abort()
			return
		}

		// Extract the token from the Authorization header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Decode the token
		claims, err := utils.DecodeToken(tokenString)
		if err != nil {
			fmt.Println(err.Error())
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = 20002
			default:
				code = 20001
			}
		}

		if code != 200 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "Token is invalid",
				"data": data,
			})

			ctx.Abort()
			return
		}

		// Set token claims in the context for use in other handlers
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
