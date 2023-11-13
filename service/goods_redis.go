package service

import (
	"blue/entity"
	"blue/global"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"math/rand"
	"time"
)

// GoFeed 确保redis中有feed项
func GoFeed() error {
	var goods *[]entity.Goods
	//查询数据库
	global.DB.Order("created_at desc").Limit(global.MaxNumGoods).Find(&goods)
	//加入缓存
	AddGoodsByGoodsIdFromCacheToRedis(goods)
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

func AddGoodsByGoodsIdFromCacheToRedis(goodsList *[]entity.Goods) error {
	for _, goods := range *goodsList {
		// 定义 key
		goodsRedis := fmt.Sprintf(entity.GoodsPattern, goods.GoodsId)
		// 使用 pipeline
		_, err := global.REDIS.TxPipelined(global.CONTEXT, func(pipe redis.Pipeliner) error {
			pipe.HSet(global.CONTEXT, goodsRedis, "goods_id", goods.GoodsId)
			pipe.HSet(global.CONTEXT, goodsRedis, "picture_url", goods.PictureUrl)
			pipe.HSet(global.CONTEXT, goodsRedis, "description", goods.Description)
			pipe.HSet(global.CONTEXT, goodsRedis, "title", goods.Title)
			pipe.HSet(global.CONTEXT, goodsRedis, "user_id", goods.UserId)
			pipe.HSet(global.CONTEXT, goodsRedis, "created_at", goods.CreatedAt)
			// 设置过期时间
			pipe.Expire(global.CONTEXT, goodsRedis, global.GoodsExpire+time.Duration(rand.Float64()*global.ExpireTimeJitter.Seconds())*time.Second)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetGoodsByGoodsIdFromRedis(goodsId int64) (*entity.Goods, error) {
	// 定义 key
	goodsRedis := fmt.Sprintf(entity.GoodsPattern, goodsId)

	var goods entity.Goods
	if result := global.REDIS.Exists(global.CONTEXT, goodsRedis).Val(); result <= 0 {
		return nil, errors.New("not found in cache")
	}
	// 使用 pipeline
	cmds, err := global.REDIS.TxPipelined(global.CONTEXT, func(pipe redis.Pipeliner) error {
		pipe.HGetAll(global.CONTEXT, goodsRedis)
		// 设置过期时间
		pipe.Expire(global.CONTEXT, goodsRedis, global.UserInfoExpire+time.Duration(rand.Float64()*global.ExpireTimeJitter.Seconds())*time.Second)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err = cmds[0].(*redis.StringStringMapCmd).Scan(&goods); err != nil {
		return nil, err
	}
	return &goods, nil
}

func GetGoodsListFromRedis() (*[]entity.Goods, error) {
	// 定义 key
	goodsRedis := "Goods:*"
	var goods []entity.Goods
	result, err := global.REDIS.Keys(global.CONTEXT, goodsRedis).Result()
	if err != nil {
		return nil, err
	}
	goods = make([]entity.Goods, len(result))
	for i, str := range result {
		cmds, err := global.REDIS.TxPipelined(global.CONTEXT, func(pipe redis.Pipeliner) error {
			pipe.HGetAll(global.CONTEXT, str)
			// 设置过期时间
			pipe.Expire(global.CONTEXT, str, global.UserInfoExpire+time.Duration(rand.Float64()*global.ExpireTimeJitter.Seconds())*time.Second)
			return nil
		})
		if err != nil {
			return nil, err
		}
		if err = cmds[0].(*redis.StringStringMapCmd).Scan(&goods[i]); err != nil {
			return nil, err
		}
	}
	return &goods, nil
}
