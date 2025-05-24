package handler

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

// SQLHandler ...
type SQLHandler struct {
	DB  *gorm.DB
	Err error
}

var dbConn *SQLHandler

func DBOpen() {
	dbConn = NewSQLHandler()
}

func DBClose() {
	sqlDB, _ := dbConn.DB.DB()
	sqlDB.Close()
}

// sql handler
func NewSQLHandler() *SQLHandler {
	// sql への接続
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	// 接続情報を表示
	// NOTE: これが表示されないと source 忘れ
	fmt.Println(user, password, host, port)
	var db *gorm.DB
	var err error
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?parseTime=true&loc=Asia%2FTokyo"
	// dsn := "docker:docker@tcp(host.docker.internal:3306)/test_database?parseTime=true&loc=Asia%2FTokyo"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	// コネクションプールの最大接続数を設定
	sqlDB.SetMaxIdleConns(10)
	//接続の最大数を設定
	sqlDB.SetMaxOpenConns(100)
	//接続の再利用が可能な時間を設定
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	sqlHandler := new(SQLHandler)
	db.Logger.LogMode(4)
	sqlHandler.DB = db

	return sqlHandler
}

// get connection
func GetDBConn() *SQLHandler {
	return dbConn
}

// bigin transaction
func BeginTransaction() *gorm.DB {
	dbConn.DB = dbConn.DB.Begin()
	return dbConn.DB
}

// rollback
func RollBack() {
	dbConn.DB.Rollback()
}
