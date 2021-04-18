package migration

import (
	"myapp/config"
	"myapp/graph/model"
)

//MigrateTable Migrate Table
func MigrateTable() {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	objectModel := []interface{}{&model.User{}, &model.Board{}, &model.List{}, &model.ListItem{}, &model.Team{}, &model.TeamHasMember{}}

	for _, val := range objectModel {
		db.AutoMigrate(val)
	}

}
