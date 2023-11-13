package initialize

import (
	"blue/entity"
	"blue/global"
	"blue/service"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() error {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/blue?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	//自动迁移
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Goods{})
	db.AutoMigrate(&entity.Comment{})
	db.AutoMigrate(&entity.Follow{})
	db.AutoMigrate(&entity.Like{})
	fmt.Println("db init")
	//u1 := User{Id: 1, Name: "张三", Gender: "男", Hobby: "学习"}
	//db.Create(&u1) //创建
	global.DB = db
	return nil
}
func InitRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.RedisHost, global.RedisPort),
		Password: global.RedisPwd,
		DB:       global.RedisDb,
		PoolSize: global.RedisPoolSize,
	})
	// 检查 Redis 连通性
	if _, err := rdb.Ping(global.CONTEXT).Result(); err != nil {
		panic(err.Error())
	}
	global.REDIS = rdb
	//先查询 feed,导入缓存
	if err := service.GoFeed(); err != nil {
		panic(err.Error())
	}
	return nil
}
