package util

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (int, int) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10 // default limit
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1 // default page
	}

	return limit, page
}

func GetPayload(c *gin.Context) (*JwtCustomClaims, error) {
	claims, exists := c.Get("payload")
	if !exists {
		return nil, fmt.Errorf("Claims not found in context")
	}
	payload, ok := claims.(*JwtCustomClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}
	return payload, nil
}
