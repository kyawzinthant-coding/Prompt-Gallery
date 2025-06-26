package database

import (
	"PromptGallery/internal/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(dtabaseURL string, environment string) error {
	var err error

	config := &gorm.Config{
		Logger:                                   getLoggerConfig(environment),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	DB, err = gorm.Open(postgres.Open(dtabaseURL), config)

	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("❌ Failed to get SQL DB instance: %v", err)
		return err
	}

	log.Println("✅ Database connected successfully")

	sqlDB.SetMaxIdleConns(10)  // Maximum idle connections
	sqlDB.SetMaxOpenConns(100) // Maximum open connections
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Printf("❌ Database migration failed: %v", err)
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil

}

// in dev mode -> detailed logs , in production -> only errors
func getLoggerConfig(environment string) logger.Interface {
	if environment == "production" {
		return logger.Default.LogMode(logger.Error)
	}
	return logger.Default.LogMode(logger.Info)
}

func GetDb() *gorm.DB {
	return DB
}

func autoMigrate() error {
	log.Println("🔄 Running database migrations...")

	return DB.AutoMigrate(
		&models.Prompt{},
		&models.User{},
		&models.PromptRequest{},
	)
}

func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db: %w", err)
	}
	return sqlDB.Close()
}
