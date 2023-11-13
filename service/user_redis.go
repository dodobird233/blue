package service

import (
	"blue/entity"
	"blue/global"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

func GetUserInfoByUserIDFromRedis(userID int64) (*entity.UserData, error) {
	// 定义 key
	userRedis := fmt.Sprintf(entity.UserDataPattern, userID)

	var userdata entity.UserData
	if result := global.REDIS.Exists(global.CONTEXT, userRedis).Val(); result <= 0 {
		return nil, errors.New("not found in cache")
	}
	// 使用 pipeline
	cmds, err := global.REDIS.TxPipelined(global.CONTEXT, func(pipe redis.Pipeliner) error {
		pipe.HGetAll(global.CONTEXT, userRedis)
		// 设置过期时间
		pipe.Expire(global.CONTEXT, userRedis, global.UserInfoExpire+time.Duration(rand.Float64()*global.ExpireTimeJitter.Seconds())*time.Second)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err = cmds[0].(*redis.StringStringMapCmd).Scan(&userdata); err != nil {
		return nil, err
	}
	return &userdata, nil
}

func AddUserInfoByUserIDFromCacheToRedis(userData *entity.UserData) error {
	// 定义 key
	userRedis := fmt.Sprintf(entity.UserDataPattern, userData.UserId)

	// 使用 pipeline
	_, err := global.REDIS.TxPipelined(global.CONTEXT, func(pipe redis.Pipeliner) error {
		pipe.HSet(global.CONTEXT, userRedis, "user_id", userData.UserId)
		pipe.HSet(global.CONTEXT, userRedis, "name", userData.Name)
		pipe.HSet(global.CONTEXT, userRedis, "follow_count", userData.FollowCount)
		pipe.HSet(global.CONTEXT, userRedis, "follower_count", userData.FollowerCount)
		pipe.HSet(global.CONTEXT, userRedis, "avatar", userData.Avatar)
		pipe.HSet(global.CONTEXT, userRedis, "background_image", userData.BackgroundImage)
		pipe.HSet(global.CONTEXT, userRedis, "signature", userData.Signature)
		pipe.HSet(global.CONTEXT, userRedis, "total_favorited", userData.TotalFavorited)
		pipe.HSet(global.CONTEXT, userRedis, "work_count", userData.WorkCount)
		pipe.HSet(global.CONTEXT, userRedis, "favorite_count", userData.FavoriteCount)
		// 设置过期时间
		pipe.Expire(global.CONTEXT, userRedis, global.UserInfoExpire+time.Duration(rand.Float64()*global.ExpireTimeJitter.Seconds())*time.Second)
		return nil
	})
	return err
}
