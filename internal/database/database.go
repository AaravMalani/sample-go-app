package database

import (
	"fmt"

	"github.com/AaravMalani/cvwo-forum-backend/internal/env"
	"github.com/AaravMalani/cvwo-forum-backend/query"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DSNFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
)

func Setup() error {

	dsn := fmt.Sprintf(
		DSNFormat, env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort, env.DBSsl, env.DBTz,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	query.SetDefault(conn)
	return nil
}
