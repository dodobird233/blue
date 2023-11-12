package entity

type Comment struct {
	CommentId  int64  `gorm:"column:comment_id;primary_key;NOT NULL" redis:"-"`
	UserId     int64  `gorm:"column:user_id;NOT NULL" redis:"user_id"`
	GoodsId    int64  `gorm:"column:goods_id;NOT NULL" redis:"goods_id"`
	Content    string `gorm:"column:comment_text;NOT NULL;type:varchar(300)" redis:"content"`
	CreateDate string `gorm:"column:create_date;NOT NULL;type:varchar(100)" redis:"-"`
}
