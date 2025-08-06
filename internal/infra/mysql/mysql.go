package mysql

import (
	"fmt"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"github.com/mqqff/savebite-be/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.AppEnv.DBUser, env.AppEnv.DBPass, env.AppEnv.DBHost, env.AppEnv.DBPort, env.AppEnv.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[MySQL][NewConn] failed to connect to database")
		return nil, err
	}

	return db, nil
}
