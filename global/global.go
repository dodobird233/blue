package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

var ( // 全局变量
	DB                  *gorm.DB      // 数据库接口
	REDIS               *redis.Client // Redis 缓存接口
	RedisHost           = "localhost"
	RedisPort           = 6379
	RedisUserName       = "root"
	RedisPwd            = "foobared"
	RedisPoolSize       = 1
	RedisDb             = 1
	FavoriteExpire      = 8 * time.Minute
	GoodsCommentsExpire = 9 * time.Minute
	CommentExpire       = 10 * time.Minute
	FollowExpire        = 11 * time.Minute
	UserInfoExpire      = 12 * time.Minute
	GoodsExpire         = 13 * time.Minute
	PublishExpire       = 14 * time.Minute
	EmptyExpire         = 15 * time.Minute
	ExpireTimeJitter    = 5 * time.Minute
	CONTEXT             = context.Background() // 上下文信息
	MaxNumGoods         = 30                   // 一次最大搜查商品量
	PATH_GOODS          = "./public/goods/"    // 商品保存相对路径
	HEAD_URL            = "http://"
	GOODS_URL           = "/static/goods/"
	ACCESS_KEY_ID       = "LTAI5tFPj7zXokBPyWyGFGx6"
	ACCESS_SEC          = "uWXWgfnzNZmPY8uHFFAYySewZIMtPQ"
)
