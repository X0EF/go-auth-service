package util

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ostheperson/go-auth-service/internal/helper"
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
		return nil, fmt.Errorf(helper.ErrFailParsePayload)
	}
	payload, ok := claims.(*JwtCustomClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}
	return payload, nil
}

const (
	charset     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	hyphenIndex = 2
)

func ParseStartAndEndTime(s, e string) (time.Time, time.Time, error) {
	jsDateFormat := "2006-01-02T15:04:05.000Z"
	start, err := time.Parse(jsDateFormat, s)
	if err != nil {
		fmt.Println(err.Error())
		return time.Time{}, time.Time{}, errors.New("Invalid start time format")
	}
	end, err := time.Parse(jsDateFormat, e)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Invalid end time format")
	}
	return start, end, nil
}

func ValidateStartAndEndTime(start, end time.Time) error {
	if start.After(end) {
		return errors.New("Start time must be before end time")
	}
	if end.Before(time.Now()) || start.Before(time.Now()) {
		return errors.New("Must Start/End reservation in future date")
	}
	// if start.Minute() != 0 || start.Second() != 0 || end.Minute() != 0 || end.Second() != 0 {
	// 	return time.Time{}, time.Time{}, errors.New("Start/End time must be on the hour")
	// }
	return nil
}

func GenerateCode(length uint) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
