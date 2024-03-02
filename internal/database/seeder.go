package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func SeedDB() {
	var env domain.Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	db := New(&env)

	for _, seed := range All() {
		log.Printf("running %v seed", seed.Name)
		if err := seed.Run(db.GetClient()); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}

func CreateUser(
	db *gorm.DB,
	firstname string,
	lastname string,
	email string,
	username string,
	password string,
	role domain.Role,
) error {
	hash, e := bcrypt.GenerateFromPassword([]byte(password), 10)
	if e != nil {
		return fmt.Errorf(helper.ErrFailHash)
	}
	user := &domain.User{
		Firstname:       firstname,
		Lastname:        lastname,
		Username:        username,
		Email:           email,
		Password:        string(hash),
		Role:            role,
		VerifiedAt:      time.Now(),
		IsEmailVerified: true,
	}
	existingUser := domain.User{}
	err := db.Where("email = ? OR username = ?", email, username).First(&existingUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(&user).Error
	}
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Email = user.Email
	existingUser.Username = user.Username
	existingUser.Password = user.Password
	existingUser.Role = user.Role
	// existingUser.IsVerified = true
	// existingUser.VerifiedAt = time.Now()
	return db.Save(&existingUser).Error
}

func All() []Seed {
	return []Seed{
		{
			Name: "CreateDefaultAdmin",
			Run: func(db *gorm.DB) error {
				return CreateUser(
					db,
					"default",
					"admin",
					os.Getenv("DEFAULT_ADMIN_EMAIL"),
					"a",
					os.Getenv("DEFAULT_ADMIN_PASSWORD"),
					domain.AdminRole,
				)
			},
		},
		{
			Name: "CreateDefaultUser",
			Run: func(db *gorm.DB) error {
				return CreateUser(
					db,
					"default",
					"user",
					os.Getenv("DEFAULT_USER_EMAIL"),
					"b",
					os.Getenv("DEFAULT_USER_PASSWORD"),
					domain.UserRole,
				)
			},
		},
	}
}
