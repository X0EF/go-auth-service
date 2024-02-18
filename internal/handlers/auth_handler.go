package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
	"github.com/ostheperson/go-auth-service/internal/util"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// User struct
}

type AuthHandler struct {
	*domain.Server
}

func NewAuthHandler(s *domain.Server) *AuthHandler {
	return &AuthHandler{Server: s}
}

func (s *AuthHandler) SignUp(c *gin.Context) {
	var newUser struct {
		Email    string
		Username string
		Password string
	}

	if c.Bind(&newUser) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": helper.ErrFailedReadBody,
		})
		return
	}
	user := domain.Users{}
	if err := s.Db.GetClient().Where("email = ? OR username = ?", newUser.Email, newUser.Username).First(&user).Error; err != gorm.ErrRecordNotFound {
		if user.Username == newUser.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrExistingUsername})
			return
		}
		if user.Email == newUser.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrExistingEmail})
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": helper.ErrFailHash,
		})
		return
	}

	user = domain.Users{
		Email:    newUser.Email,
		Password: string(hash),
		Role:     domain.UserRole,
		Username: newUser.Username,
	}
	result := s.Db.GetClient().Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": helper.ErrFailCreateUser,
		})
		return
	}

	// Respond
	c.JSON(http.StatusCreated, domain.Response{
		Message: helper.Success,
		Data:    &user,
	})
}

func (s *AuthHandler) SignIn(c *gin.Context) {
	var details struct {
		Email    string
		Username string
		Password string
	}

	err := c.ShouldBind(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := domain.Users{}
	if err := s.Db.GetClient().Where("email = ? OR username = ?", details.Email, details.Username).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			if user.Username == details.Username {
				c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrNoExistingUsername})
				return
			}
			if user.Email == details.Email {
				c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrNoExistingEmail})
				return
			}
		}
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(details.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	accessToken, err := util.CreateAccessToken(
		&user,
		s.Env.AccessTokenSecret,
		s.Env.AccessTokenExpiryHour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err := util.CreateRefreshToken(
		&user,
		s.Env.RefreshTokenSecret,
		s.Env.RefreshTokenExpiryHour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusCreated, domain.Response{
		Message: helper.Success,
		Data:    loginResponse,
	})
}

func (s *AuthHandler) SignInAdmin(c *gin.Context) {
	var details struct {
		Email    string
		Username string
		Password string
	}

	err := c.ShouldBind(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := domain.Users{}
	if err := s.Db.GetClient().
		Where("email = ? OR username = ? AND role = ?", details.Email, details.Username, domain.AdminRole).
		First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			if user.Username == details.Username {
				c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrNoExistingUsername})
				return
			}
			if user.Email == details.Email {
				c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrNoExistingEmail})
				return
			}
		}
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(details.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	accessToken, err := util.CreateAccessToken(
		&user,
		s.Env.AccessTokenSecret,
		s.Env.AccessTokenExpiryHour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err := util.CreateRefreshToken(
		&user,
		s.Env.RefreshTokenSecret,
		s.Env.RefreshTokenExpiryHour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusCreated, domain.Response{
		Message: helper.Success,
		Data:    loginResponse,
	})
}
