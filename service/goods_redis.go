package service

import (
	"blue/entity"
	"blue/global"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"math/rand"
	"time"
)

// GoFeed 确保redis中有feed项
func GoFeed() error {
	n, err := global.REDIS.Exists(global.CONTEXT, "feed").Result()
	if err != nil {
		return err
	}
	if n <= 0 {
		// "feed"不存在
		var allGoods []entity.Goods
		if err := global.DB.Find(&allGoods).Error; err != nil {
			return err
		}
		if len(allGoods) == 0 {
			return nil
		}
		var listZ = make([]*redis.Z, 0, len(allGoods))
		for _, goods := range allGoods {
			listZ = append(listZ, &redis.Z{Score: float64(goods.CreatedAt.UnixMilli()) / 1000, Member: goods.GoodsId})
		}
		return global.REDIS.ZAdd(global.CONTEXT, "feed", listZ...).Err()
	}
	return nil
}

func SaveCommentsOfGoods(commentList []entity.Comment, keyCommentsOfGoods string) error {
	var listZ = make([]*redis.Z, 0, len(commentList))
	for _, comment := range commentList {
		t, _ := time.Parse(comment.CreateDate, "01-02")
		listZ = append(listZ, &redis.Z{Score: float64(t.UnixMilli()) / 1000, Member: comment.CommentId})
	}
	pipe := global.REDIS.TxPipeline()
	pipe.ZAdd(global.CONTEXT, keyCommentsOfGoods, listZ...)
	pipe.Expire(global.CONTEXT, keyCommentsOfGoods, global.GoodsCommentsExpire+time.Duration(rand.Float64()*global.ExpireTimeJitter.Seconds())*time.Second)
	_, err := pipe.Exec(global.CONTEXT)
	return err
}

func GetCommentCountOfGoods(goodsId int64) (int, error) {
	keyGoods := fmt.Sprintf(entity.GoodsPattern, goodsId)
	lua := redis.NewScript(`
				local key = KEYS[1]
				local expire_time = ARGV[1]
				if redis.call("Exists", key) > 0 then
					redis.call("Expire", key, expire_time)
					return redis.call("HGet", key, "comment_count")
				end
				return -1
			`)
	keys := []string{keyGoods}
	values := []interface{}{global.GoodsCommentsExpire.Seconds() + math.Floor(rand.Float64()*global.ExpireTimeJitter.Seconds())}
	numComments, err := lua.Run(global.CONTEXT, global.REDIS, keys, values).Int()
	if err != nil {
		return 0, err
	}
	return numComments, nil
}
