package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
	"github.com/ostheperson/go-auth-service/internal/util"
)

type UsersHandler struct {
	*domain.Server
}

func NewUsersHandler(s *domain.Server) *UsersHandler {
	return &UsersHandler{Server: s}
}

func (s *UsersHandler) GetUsers(c *gin.Context) {
	var users []domain.Users
	limit, page := util.GetPaginationParams(c)

	// TODO: Add filtering options
	if err := s.Db.GetClient().Limit(limit).Offset(page - 1).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.ErrFailGetUsers})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
		Data:    &users,
	})
}

func (s *UsersHandler) GetUser(c *gin.Context) {
	payload, err := util.GetPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.ErrFailParsePayload})
		return
	}
	id := c.Param("id")
	if fmt.Sprint(payload.ID) != id && payload.Role != domain.AdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": helper.ErrUnauthorized})
		return
	}
	user := domain.Users{}
	if err := s.Db.GetClient().First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": helper.NotFound("users")})
			return
		}
		c.JSON(500, gin.H{"error": helper.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
		Data:    &user,
	})
}

func (s *UsersHandler) UpdateUser(c *gin.Context) {
	var details struct {
		Firstname *string
		Lastname  *string
		Bio       *string
		AvatarURL *string
	}

	if c.Bind(&details) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": helper.ErrFailedReadBody,
		})
		return
	}
	id := c.Param("id")
	user := domain.Users{}
	if err := s.Db.GetClient().First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": helper.NotFound("user")})
		} else {
			c.JSON(500, gin.H{"error": helper.ErrInternalError})
		}
		return
	}
	if details.Firstname != nil {
		user.Firstname = *details.Firstname
	}
	if details.Lastname != nil {
		user.Lastname = *details.Lastname
	}
	if details.Bio != nil {
		user.Bio = *details.Bio
	}
	if details.AvatarURL != nil {
		user.AvatarURL = *details.AvatarURL
	}
	user.UpdatedAt = time.Now()
	if err := s.Db.GetClient().Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": helper.ErrInternalError})
	}
	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
		Data:    &user,
	})
}

func (s *UsersHandler) RemoveUser(c *gin.Context) {
	payload, err := util.GetPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	if fmt.Sprint(payload.ID) != id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": helper.ErrUnauthorized})
		return
	}
	if err := s.Db.GetClient().Model(&domain.Users{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		c.JSON(500, gin.H{"error": helper.ErrFailDelUser})
		return
	}
	c.JSON(200, gin.H{"message": helper.Success})
}

func (s *UsersHandler) ClearTable(c *gin.Context) {
	if err := s.Db.GetClient().Delete(&domain.Users{}).Error; err != nil {
		panic("failed to clear table")
	}
}
