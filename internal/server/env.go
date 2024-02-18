package server

import (
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/ostheperson/go-auth-service/internal/domain"
)

func NewEnv() *domain.Env {
	var env domain.Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
