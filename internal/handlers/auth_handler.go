package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/integrations"
	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
	"github.com/ostheperson/go-auth-service/internal/services"
	"github.com/ostheperson/go-auth-service/internal/util"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	*domain.Server
	mailer integrations.MailerService
	ts     services.TokenService
}

func NewAuthHandler(
	s *domain.Server,
	mailer integrations.MailerService,
	tokenService services.TokenService,
) *AuthHandler {
	return &AuthHandler{s, mailer, tokenService}
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
	user := domain.User{}
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

	user = domain.User{
		Email:    newUser.Email,
		Password: string(hash),
		Role:     domain.UserRole,
		Username: newUser.Username,
	}
	result := s.Db.GetClient().Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": helper.FailCreate("user"),
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
	user := domain.User{}
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": helper.ErrInvalidCredentials})
		return
	}
	accessToken, refreshToken, err := login(&user, s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, domain.Response{
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
	user := domain.User{}
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": helper.ErrInvalidCredentials})
		return
	}
	accessToken, refreshToken, err := login(&user, s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
		Data:    loginResponse,
	})
}

func (s *AuthHandler) ForgotPassword(c *gin.Context) {
	var details struct {
		Email string `json:"email"`
	}

	err := c.ShouldBind(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := domain.User{}
	err = s.Db.GetClient().Where("email = ?", details.Email).First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrNoExistingEmail})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	code := util.GenerateCode(s.Env.ConfirmCodeLength)
	expires := time.Now().Add(time.Duration(s.Env.ConfirmationCodeExpiryHour) * time.Hour)
	_, err = s.ts.CreateToken(domain.RESET_PASSWORD, code, user.ID, expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.ErrGenerateToken})
		return
	}

	err = s.mailer.SendMail(
		s.Env.POSTMARK_FROM_EMAIL,
		details.Email,
		"Password reset code",
		fmt.Sprintf("your reset code is %s", code),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.ErrCannotSendMail})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Message: helper.CodeSentToEmail(details.Email),
	})
	return
}

func (s *AuthHandler) ResetPassword(c *gin.Context) {
	var details struct {
		Email    string `json:"email"`
		Password string `json:"new_password"`
		Code     string `json:"code"`
	}

	err := c.ShouldBind(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := domain.Token{}
	err = s.Db.GetClient().Preload("User").Where("hash = ?", details.Code).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrInvalidCode})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if time.Now().After(token.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrExpiredCode})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(details.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": helper.ErrFailHash,
		})
		return
	}
	token.User.Password = string(hash)

	if err := s.Db.GetClient().Save(&token.User).Error; err != nil {
		c.JSON(500, gin.H{"error": helper.ErrInternalError})
		return
	}
	if err := s.Db.GetClient().Delete(&token).Error; err != nil {
		c.JSON(500, gin.H{"error": helper.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
	})
	return
}

func (s *AuthHandler) ResendVerifyEmail(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	user := domain.User{}
	err := s.Db.GetClient().Where("email = ?", email).First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrNoExistingEmail})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	s.Db.GetClient().
		Where("type = ? AND userid = ?", domain.RESET_PASSWORD, user.ID).
		Delete(&domain.User{})
	code := util.GenerateCode(s.Env.ConfirmCodeLength)
	s.L.Print(code)
	expires := time.Now().Add(time.Duration(s.Env.ConfirmationCodeExpiryHour) * time.Hour)
	_, err = s.ts.CreateToken(domain.RESET_PASSWORD, code, user.ID, expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.ErrGenerateToken})
		return
	}

	err = s.mailer.SendMail(
		s.Env.POSTMARK_FROM_EMAIL,
		email,
		"Password reset code",
		fmt.Sprintf("your reset code is %s", code),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.ErrCannotSendMail})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Message: fmt.Sprintf("Verification code sent to email %s", email),
	})
	return
}

func (s *AuthHandler) ConfirmEmail(c *gin.Context) {
	var details struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	err := c.ShouldBind(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := domain.Token{}
	err = s.Db.GetClient().Preload("User").Where("hash = ?", details.Code).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrInvalidCode})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if time.Now().After(token.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrExpiredCode})
		return
	}

	token.User.IsEmailVerified = true
	if err := s.Db.GetClient().Save(&token.User).Error; err != nil {
		c.JSON(500, gin.H{"error": helper.ErrInternalError})
		return
	}
	if err := s.Db.GetClient().Delete(&token).Error; err != nil {
		c.JSON(500, gin.H{"error": helper.ErrInternalError})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
	})
	return
}

func (s *AuthHandler) RefreshToken(c *gin.Context) {
	var details struct {
		Hash string `json:"hash"`
	}
	err := c.ShouldBind(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := domain.Token{}
	err = s.Db.GetClient().Preload("User").Where("hash = ?", details.Hash).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.ErrInvalidRefreshToken})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprint(err.Error() + "when validating refresh token")},
		)
		return
	}
	accessToken, refreshToken, err := login(&token.User, s)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprint(err.Error() + "when generating login tokens")},
		)
		return
	}
	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, domain.Response{
		Message: helper.Success,
		Data:    loginResponse,
	})
}

func login(user *domain.User, s *AuthHandler) (string, string, error) {
	// TODO: Handle devices
	accessToken, err := util.CreateAccessToken(
		user,
		s.Env.AccessTokenSecret,
		s.Env.AccessTokenExpiryHour,
	)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := util.CreateRefreshToken(
		user,
		s.Env.RefreshTokenSecret,
		s.Env.RefreshTokenExpiryHour,
	)
	if err != nil {
		return "", "", err
	}
	expires := time.Now().Add(time.Duration(s.Env.RefreshTokenExpiryHour) * time.Hour)
	_, err = s.ts.CreateToken(domain.REFRESH, refreshToken, user.ID, expires)
	if err != nil {
		return "", "", err
	}
	user.LastLoggedInAt = time.Now()
	if err := s.Db.GetClient().Save(&user).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
