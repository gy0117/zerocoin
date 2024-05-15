package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"zero-common/zerodb"
)

const (
	maxOpenConnections = 100
	maxIdleConnections = 10
)

// ConnMysql 连接数据库
func ConnMysql(dsn string) *zerodb.ZeroDB {
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect mysql, err: " + err.Error())
	}
	db, err := _db.DB()
	if err != nil {
		panic("failed to connect mysql, err: " + err.Error())
	}

	// 连接池配置
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	return &zerodb.ZeroDB{Conn: _db}
}
