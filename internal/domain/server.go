package domain

import (
	"log"

	"gorm.io/gorm"
)

type Service interface {
	Health() map[string]string
	GetClient() *gorm.DB
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Server struct {
	Port int
	Db   Service
	Env  *Env
	L    *log.Logger
}
