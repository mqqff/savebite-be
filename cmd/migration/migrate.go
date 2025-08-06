package main

import (
	"github.com/mqqff/savebite-be/internal/domain/entity"
	"github.com/mqqff/savebite-be/internal/infra/mysql"
	"github.com/mqqff/savebite-be/pkg/log"
)

func main() {
	db, err := mysql.NewConn()
	if err != nil {
		return
	}

	err = db.Migrator().DropTable(&entity.User{}, &entity.Analysis{}, &entity.Ingredient{})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Migration][main] failed to drop all tables")
		return
	}

	err = db.Migrator().CreateTable(&entity.User{}, &entity.Analysis{}, &entity.Ingredient{})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Migration][main] failed to create tables")
		return
	}

	log.Info(nil, "[Migration][main] migration success")
}
