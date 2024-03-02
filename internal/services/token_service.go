package services

import (
	"time"

	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/internal/domain"
)

type TokenService interface {
	CreateToken(
		ttype domain.TokenType,
		hash string,
		UserID uint,
		expires time.Time,
	) (*domain.Token, error)
	GetTokenByID(id uint) (*domain.Token, error)
	DeleteTokenByID(id uint) error
}
type tokenService struct {
	db *gorm.DB
}

func NewTokenService(db *gorm.DB) TokenService {
	return &tokenService{db: db}
}

func (s *tokenService) CreateToken(
	ttype domain.TokenType,
	hash string,
	UserID uint,
	expires time.Time,
) (*domain.Token, error) {
	token := &domain.Token{Hash: hash, Type: ttype, UserID: UserID, ExpiresAt: expires}
	if err := s.db.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (s *tokenService) GetTokenByID(id uint) (*domain.Token, error) {
	token := &domain.Token{}
	if err := s.db.First(token, id).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (s *tokenService) DeleteTokenByID(id uint) error {
	return s.db.Delete(&domain.Token{}, id).Error
}
