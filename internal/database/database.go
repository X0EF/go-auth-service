package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ostheperson/go-auth-service/internal/domain"
)

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second,   // Slow SQL threshold
		LogLevel:                  logger.Silent, // Log level
		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      true,          // Don't include params in the SQL log
		Colorful:                  false,         // Disable color
	},
)

type dbclient struct {
	db *gorm.DB
}

func New(env *domain.Env) domain.DBService {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DB_USERNAME,
		env.DB_PASSWORD,
		env.DB_HOST,
		env.DB_PORT,
		env.DB_DATABASE,
	)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		// NowFunc: func() time.Time {
		// 	ti, _ := time.LoadLocation(env.Timezone)
		// 	return time.Now().In(ti)
		// },
		// Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("SET TIMEZONE='UTC'")
	// postgresdb, err := db.DB()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer postgresdb.Close()
	s := &dbclient{db: db}
	return s
}

func (s *dbclient) GetClient() *gorm.DB {
	return s.db
}

func (s *dbclient) Health() map[string]string {
	sqlDB, err := s.db.DB()
	if err != nil {
		panic("failed to get underlying *sql.DB instance")
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *dbclient) AutoMigrateAll(models []interface{}) {
	for _, model := range models {
		s.db.AutoMigrate(model)
	}
}
