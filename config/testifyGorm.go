package config

import (
	"fmt"
	logger "myapp/logger"

	mysql "gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
)

//Generated By github.com/david-yappeter/GormCrudGenerator

//ConnectGorm Database Connection to Gorm V2
func ConnectGormTest() *gorm.DB {
	databaseConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", "root", "root", "localhost", "3306", "test")

	db, err := gorm.Open(mysql.Open(databaseConfig), logger.InitConfigTest())

	if err != nil {
		fmt.Println(err)
		panic("Fail To Connect Database")
	}

	return db
}
