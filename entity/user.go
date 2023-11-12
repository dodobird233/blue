package entity

type UserData struct {
	UserId          int64  `json:"id" redis:"user_id"`
	Name            string `json:"name" redis:"name"`
	FollowCount     int64  `json:"follow_count" redis:"follow_count"`
	FollowerCount   int64  `json:"follower_count" redis:"follower_count"`
	IsFollow        bool   `json:"is_follow" redis:"is_follow"`
	Avatar          string `json:"avatar" redis:"avatar"`
	BackgroundImage string `json:"background_image" redis:"background_image"`
	Signature       string `json:"signature" redis:"signature"`
	TotalFavorited  int64  `json:"total_favorited" redis:"total_favorited"`
	WorkCount       int64  `json:"work_count" redis:"work_count"`
	FavoriteCount   int64  `json:"favorite_count" redis:"favorite_count"`
}

type User struct {
	UserId          int64  `gorm:"column:user_id;primary_key;NOT NULL" redis:"user_id"`
	UserName        string `gorm:"column:user_name;type:varchar(100)" redis:"user_name"`
	Password        string `gorm:"column:password;type:varchar(100)" redis:"password"`
	Avatar          string `gorm:"column:avatar;type:varchar(100)" redis:"avatar"`
	BackgroundImage string `gorm:"column:background_image;type:varchar(100)" redis:"background_image"`
	Signature       string `gorm:"column:signature;type:varchar(200)" redis:"signature"`
}
