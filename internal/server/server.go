package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ostheperson/go-auth-service/internal/database"
	"github.com/ostheperson/go-auth-service/internal/domain"
)

func NewServer() *http.Server {
	l := log.New(os.Stdout, "auth-api ", log.LstdFlags)
	env := NewEnv()
	db := database.New(env)
	// db.GetClient().Migrator().DropTable("users")
	db.GetClient().AutoMigrate(&domain.Users{})
	NewServer := &domain.Server{
		Port: env.PORT,
		Db:   db,
		Env:  env,
		L:    l,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      RegisterRoutes(NewServer),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     l,
	}

	return server
}
