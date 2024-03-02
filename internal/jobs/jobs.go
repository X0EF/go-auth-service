package jobs

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/ostheperson/go-auth-service/internal/domain"
)

type JobService interface {
	Start()
}

type jobService struct {
	cron *cron.Cron
	db   *gorm.DB
	l    *log.Logger
}

func NewJobService(tz string, db *gorm.DB) JobService {
	l := log.New(os.Stdout, "jobs", log.LstdFlags)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	return &jobService{cron: cron.New(cron.WithLocation(loc)), l: l, db: db}
}

func (js *jobService) Start() {
	log.Println("Starting job service...")
	js.scheduleJobs()
	js.cron.Start()
}

func (js *jobService) scheduleJobs() {
	// remove expired tokens
	js.cron.AddFunc("@every 5m", func() {
		js.db.Where("expires_at < ?", time.Now()).Delete(&domain.Token{})
	})

	// TODO: add job to delete pending payment links
}
