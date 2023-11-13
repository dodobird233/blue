package entity

import "time"

type Goods struct {
	GoodsId     int64     `gorm:"column:goods_id;primary_key;NOT NULL" redis:"goods_id"`
	PictureUrl  string    `gorm:"column:picture_url;type:varchar(500)" redis:"picture_url"`
	Description string    `gorm:"column:description;type:varchar(500)" redis:"description"`
	Title       string    `gorm:"column:title;type:varchar(100)" redis:"title"`
	UserId      int64     `gorm:"column:user_id;NOT NULL" redis:"user_id"`
	CreatedAt   time.Time `gorm:"column:created_at;index" redis:"-"`
}

type GoodsLikeCnt struct {
	GoodsId int64 `gorm:"column:goods_id;primary_key;NOT NULL"`
	LikeCnt int64 `gorm:"column:like_cnt"`
}

type GoodsCommentCnt struct {
	GoodsId    int64 `gorm:"column:goods_id;primary_key;NOT NULL"`
	CommentCnt int64 `gorm:"column:comment_cnt"`
}
