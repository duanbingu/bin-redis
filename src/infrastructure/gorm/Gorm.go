package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var GormDB *gorm.DB
func init(){
	dsn:="root:root@tcp(localhost:3306)/jichuban?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB,err:=db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetMaxOpenConns(10)
	GormDB =db
}