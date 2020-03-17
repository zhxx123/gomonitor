package db

import (
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
)

// DB 数据库连接
var DB *gorm.DB

// RedisPool Redis连接池
var RedisPool *redis.Pool

// MongoDB 数据库连接
var MongoDB *mgo.Database

// 初始化数据库
func initDB() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Printf("No error should happen when connecting to  database, but got err=%s", err.Error())
		os.Exit(-1)
	}
	if config.ServerConfig.Env == model.DevelopmentMode {
		db.LogMode(true) //开启日志模式
	}
	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)
	DB = db

	// 初始化数据库表
	if config.DBConfig.AutoMigrated == true {
		if err = DB.AutoMigrate(model.Models...).Error; nil != err {
			fmt.Printf("auto migrate tables failed: %s", err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println("init DB", config.DBConfig.Dialect)
}

// 数据库关闭
func closeDB() {
	if err := DB.Close(); nil != err {
		fmt.Printf("[IrisDB] Close db err %s", err.Error())
	}
	fmt.Println("[IrisDB] Close")
}

func initRedis() {
	RedisPool = &redis.Pool{
		MaxIdle:     config.RedisConfig.MaxIdle,
		MaxActive:   config.RedisConfig.MaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisConfig.URL, redis.DialPassword(config.RedisConfig.Password))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func closeRedis() {
	if err := RedisPool.Close(); err != nil {
		fmt.Printf("[RedisDB] Close redis db err %s", err.Error())
	}
	fmt.Println("[RedisDB] Close")
}

/*
 * mgo文档 http://labix.org/mgo
 * https://godoc.org/gopkg.in/mgo.v2
 * https://godoc.org/gopkg.in/mgo.v2/bson
 * https://godoc.org/gopkg.in/mgo.v2/txn
 */
func initMongo() {
	if config.MongoConfig.URL == "" {
		return
	}
	session, err := mgo.Dial(config.MongoConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	MongoDB = session.DB(config.MongoConfig.Database)
}

// 初始化
func InitDB() {
	initDB()
	initRedis()
	initMongo()
}

// 关闭
func CloseAllDB() {
	closeDB()
	closeRedis()
}
