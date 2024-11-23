package db

import (
	"common/zerodb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

// ConnMysql2 支持读写分离
func ConnMysql2(dsn string, slaveDsn string) *zerodb.ZeroDB {
	masterDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect to  master DB, err: " + err.Error())
	}

	// 从库连接
	slaveDB, err := gorm.Open(mysql.Open(slaveDsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect to slave DB, err: " + err.Error())
	}

	db, _ := masterDB.DB()
	// 连接池配置
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	// 配置主从分离
	masterDB.Set("gorm:read_default", slaveDB)
	return &zerodb.ZeroDB{Conn: masterDB}
}
