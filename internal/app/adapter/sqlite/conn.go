package sqlite

import (
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connection gets connection of sqlite database
func Connection() (db *gorm.DB) {
	dsn := viper.GetString("database.sqlite.path")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// can be used for debugging
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
