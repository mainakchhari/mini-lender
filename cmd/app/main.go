package main

import (
	"github.com/mainakchhari/mini-lender/internal/app/adapter"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite/model"
	"github.com/spf13/viper"
)

func main() {
	// Setup Default Router
	r := adapter.Router()

	// Automigrate schemas
	viper.SetDefault("database.sqlite.path", "mini-lender.db")
	db := sqlite.Connection()
	db.AutoMigrate(&model.User{}, &model.Loan{}, &model.Payment{})

	// for local development
	r.Run(":8080")
}
