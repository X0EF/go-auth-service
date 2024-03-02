package domain

import (
	"log"

	"gorm.io/gorm"
)

type DBService interface {
	Health() map[string]string
	GetClient() *gorm.DB
	AutoMigrateAll(models []interface{})
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Server struct {
	Port int
	Db   DBService
	Env  *Env
	L    *log.Logger
}
