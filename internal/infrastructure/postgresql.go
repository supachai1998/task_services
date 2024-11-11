package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/supachai1998/task_services/internal/configs"
)

func NewPostgreSQL(config *configs.DatabaseConfig) (*gorm.DB, error) {
	// Set up the database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		configs.AppConfig.Database.Host,
		configs.AppConfig.Database.User,
		configs.AppConfig.Database.Password,
		configs.AppConfig.Database.DbName,
		configs.AppConfig.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: GormLogger{logger.Default.LogMode(logger.Info)},
	})
	if err != nil {
		panic("failed to connect to database")
	}

	return db, nil
}

type GormLogger struct {
	logger.Interface
}
