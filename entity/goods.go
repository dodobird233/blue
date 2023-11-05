package entity

import "time"

type Goods struct {
	GoodsId     int64     `gorm:"column:goods_id;primary_key;NOT NULL"`
	PictureUrl  string    `gorm:"column:picture_url;type:varchar(500)"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	Title       string    `gorm:"column:title;type:varchar(100)"`
	UserId      int64     `gorm:"column:user_id;NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_at;index"`
}

type GoodsLikeCnt struct {
	GoodsId int64 `gorm:"column:goods_id;primary_key;NOT NULL"`
	LikeCnt int64 `gorm:"column:like_cnt"`
}

type GoodsCommentCnt struct {
	GoodsId    int64 `gorm:"column:goods_id;primary_key;NOT NULL"`
	CommentCnt int64 `gorm:"column:comment_cnt"`
}
