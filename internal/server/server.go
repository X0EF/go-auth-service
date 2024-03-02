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
	l := log.New(os.Stdout, "autoparts-api", log.LstdFlags)
	env := NewEnv(".env")
	db := database.New(env)
	db.AutoMigrateAll(domain.GetModels())
	// db.GetClient().Migrator().DropTble("users")
	// db.GetClient().Migrator().DropTable("listings")
	// db.GetClient().Migrator().DropTable("reservations")
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

	// jobservice := jobs.NewJobService(env.Timezone, db.GetClient())
	// jobservice.Start()

	return server
}
