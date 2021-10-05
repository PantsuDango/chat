package db

import (
	model2 "chat/src/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var exeDB *gorm.DB

// 连接数据库
func OpenDB(ConfigYaml model2.ConfigYaml) {

	ip := ConfigYaml.Mysql.IP
	port := ConfigYaml.Mysql.Port
	user := ConfigYaml.Mysql.User
	password := ConfigYaml.Mysql.Password
	database := ConfigYaml.Mysql.Database

	connectInfo := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`,
		user, password, ip, port, database)

	var err error
	exeDB, err = gorm.Open("mysql", connectInfo)
	if err != nil {
		log.Fatalf("Connect db fail: %s", err)
	}

}

// 关闭数据库
func CloseDB() {
	if exeDB != nil {
		exeDB.Close()
	}
}
