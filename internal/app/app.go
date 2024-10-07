package app

import (
	"StudentManager/internal/config"
	"StudentManager/pkg/database/postgres"
)

func Run() {
	cfg := config.Init()
	db, _ := postgres.New(cfg.Database)
	_ = db
}
