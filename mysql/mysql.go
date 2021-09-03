package mysql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	_db   *gorm.DB
	_conf *Config
)

type Config struct {
	ConnString      string
	ConnMaxLifetime int64 //ConnMaxLifetime 最大连接时间，单位：小时
	MaxIdleConns    int
	MaxOpenConns    int
}

//Init mysql初始化
func Init(conf *Config) {
	if conf != nil {
		_conf = conf
	}
}

//New 创建实例
func New() (*gorm.DB, error) {
	if _db != nil {
		return _db, nil
	}

	_db, err := gorm.Open(mysql.Open(_conf.ConnString), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	_db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(_conf.ConnString)},
			Replicas: []gorm.Dialector{mysql.Open(_conf.ConnString)},
			Policy:   dbresolver.RandomPolicy{},
		}).SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(time.Duration(_conf.ConnMaxLifetime) * time.Hour).
			SetMaxIdleConns(_conf.MaxIdleConns).
			SetMaxOpenConns(_conf.MaxOpenConns))
	return _db, nil
}
