package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
	"github.com/ostheperson/go-auth-service/internal/util"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "payload"
)

func AuthMiddleware(secret string, roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		JwtAuthMiddleware(secret)(c)
		RoleMiddleware(roles...)(c)
	}
}

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)

		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authorization header is not provided"},
			)
			return
		}

		fields := strings.Split(authHeader, " ")
		if len(fields) < 2 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid authorization header format"},
			)
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Sprintf("unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		accessToken := fields[1]
		payload, err := util.VerifyAndExtract(accessToken, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

func RoleMiddleware(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(authorizationPayloadKey)
		if !exists {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": helper.ErrFailParseRole},
			)
			return
		}
		payload, ok := claims.(*util.JwtCustomClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse claims"})
			return
		}
		userRole := payload.Role

		for _, role := range roles {
			if role == userRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
