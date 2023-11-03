package global

import (
	"gorm.io/gorm"
)

var ( // 全局变量
	DB          *gorm.DB            // 数据库接口
	MaxNumGoods = 30                // 一次最大搜查商品量
	PATH_GOODS  = "./public/goods/" // 商品保存相对路径
	HEAD_URL    = "http://"
	GOODS_URL   = "/static/goods/"
)
