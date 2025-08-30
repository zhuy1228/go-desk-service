package libs

import (
	"fmt"
	"go-desk-service/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	appConfig, dberr := config.LoadConfig()
	if dberr != nil {
		fmt.Println("读取配置失败")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		appConfig.Database.Host,
		appConfig.Database.Username,
		appConfig.Database.Password,
		appConfig.Database.Database,
		appConfig.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established")
	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("database connection is not initialized")
	}
	return DB
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
