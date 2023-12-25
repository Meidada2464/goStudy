package SyncOAData

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDataBase() (db *gorm.DB, err error) {
	dsn := "root:aresshield@tcp(172.18.156.144:31045)/qos_portal?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
