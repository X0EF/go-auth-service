package services

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/integrations"
	"github.com/ostheperson/go-auth-service/internal/domain"
)

type VerifyService interface {
	VerifyEmail(email string) error
}

type verifyService struct {
	e      *domain.Env
	db     *gorm.DB
	mailer integrations.MailerService
}

func (v *verifyService) VerifyEmail(email, code string) error {
	err := v.mailer.SendMail(
		v.e.POSTMARK_FROM_EMAIL,
		email,
		"Verify account",
		fmt.Sprintf("your verification code is %s", code),
	)
	if err != nil {
		return err
	}
	return nil
}
