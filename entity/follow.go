package entity

type Follow struct {
	FollowId     int64 `gorm:"column:follow_id;primary_key;NOT NULL"`
	UserId       int64 `gorm:"column:user_id;NOT NULL;index:follow_user_id,priority:2;index:follow_follower_id"`
	FollowUserId int64 `gorm:"column:follow_userid;NOT NULL;index:follow_user_id,priority:1"`
}
