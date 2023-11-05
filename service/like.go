package service

import (
	"blue/entity"
	"blue/global"
	"blue/util"
	"errors"
	"sort"
)

// 获取每个商品的点赞数量
func QueryLikeCountListByGoodsIdList(goodsIdList *[]int64) ([]entity.GoodsLikeCnt, error) {
	var getLikeCountList []entity.GoodsLikeCnt
	result := global.DB.Model(&entity.Like{}).Select("goods_id", "count(goods_id) as like_cnt").Where(map[string]interface{}{"goods_id": *goodsIdList}).Group("goods_id").Find(&getLikeCountList)
	if result.Error != nil {
		err := errors.New("likesList query failed")
		return nil, err
	}
	// 找数据找齐了
	if len(*goodsIdList) == len(getLikeCountList) {
		return getLikeCountList, nil
	}
	// 数据不全，误差部分补全为0
	var likeCountList []entity.GoodsLikeCnt
	likeCountList = make([]entity.GoodsLikeCnt, len(*goodsIdList))
	for i, goodsId := range *goodsIdList {
		likeCountList[i].GoodsId = goodsId
		likeCountList[i].LikeCnt = FindGoodsIdFromGoodsLikeCntList(goodsId, &getLikeCountList)
	}
	return likeCountList, nil
}

// 使用用户id查询其点赞商品的id列表
func QueryLikeGoodsIdListByUserId(userId int64) (likeList []int64, err error) {
	result := global.DB.Model(&entity.Like{}).Select("goods_id").Where("user_id=?", userId).Find(&likeList)
	if result.Error != nil {
		return nil, err
	}
	return
}

// 根据用户id以及给定商品id列表返回点赞列表情况
func ParseLikeGoodsListByUserIdFormGoodsId(userId int64, goodsIdList *[]int64) (isFavoriteList []bool, err error) {
	var likeList []int64
	likeList, err = QueryLikeGoodsIdListByUserId(userId)
	if err != nil {
		return nil, err
	}
	sort.Slice(likeList, func(i, j int) bool { return likeList[i] < likeList[j] })
	isFavoriteList = make([]bool, len(*goodsIdList))
	for i, goodsId := range *goodsIdList {
		isFavoriteList[i] = FindInt64(goodsId, &likeList)
	}
	return
}

// 根据商品id查询商品对象
func QueryGoodsListByGoodsIdList(goodsIdList *[]int64) (goodsList []entity.Goods, err error) {
	result := global.DB.Model(&entity.Goods{}).Where("goods_id in ?", *goodsIdList).Find(&goodsList)
	if result.Error != nil {
		return nil, err
	}
	return
}

// 点赞操作和取消赞操作
func GiveOrCancelLike(userId int64, goodsId int64, actionType int32) (err error) {
	var likeList []entity.Like
	result := global.DB.Model(&entity.Like{}).Where("user_id=? and goods_id=?", userId, goodsId).Find(&likeList)
	if result.Error != nil {
		return
	}
	//查询到有点赞记录
	if likeList != nil && len(likeList) > 0 {
		//已经点赞过
		if actionType == 1 {
			return
		}
		//取消点赞
		var cancelLike entity.Like
		cancelLike.LikeId = likeList[0].LikeId
		result = global.DB.Model(&entity.Like{}).Delete(&cancelLike)
		if result.Error != nil {
			return err
		}
		return
	}
	//无点赞记录
	//取消点赞
	if actionType == 2 {
		return
	}
	//进行点赞
	var giveLike entity.Like
	giveLike.LikeId = util.GetNextId()
	giveLike.UserId = userId
	giveLike.GoodsId = goodsId
	if global.DB.Model(&entity.Like{}).Create(&giveLike).Error != nil {
		return err
	}
	return
}

// 根据id查询点赞商品列表
func GetLikeGoodsListByUserId(userId int64, currentId int64) (goods []entity.GoodsResponse, err error) {
	//查询当前用户的点赞的商品id列表
	likeGoodsIdList, err := QueryLikeGoodsIdListByUserId(currentId)
	if err != nil {
		return nil, err
	}
	//根据商品id列表查询商品对象
	likeGoodsList, err := QueryGoodsListByGoodsIdList(&likeGoodsIdList)
	if err != nil {
		return nil, err
	}
	//根据商品id列表查询点赞数量
	likeCountList, err := QueryLikeCountListByGoodsIdList(&likeGoodsIdList)
	if err != nil {
		return nil, err
	}
	//防止数量为0,预先使用map记录
	likeCountListMap := map[int64]int64{}
	for _, likeCount := range likeCountList {
		likeCountListMap[likeCount.GoodsId] = likeCount.LikeCnt
	}
	//根据商品id列表查询评论数量
	commentCountList, err := QueryCommentCountListByGoodsIdList(&likeGoodsIdList)
	if err != nil {
		return nil, err
	}
	commentCountListMap := map[int64]int64{}
	for _, likeCount := range commentCountList {
		commentCountListMap[likeCount.GoodsId] = likeCount.CommentCnt
	}
	goods = make([]entity.GoodsResponse, len(likeGoodsList))
	for i, goodsItem := range likeGoodsList {
		goods[i].Id = goodsItem.GoodsId
		goods[i].Author, err = UserInfoByUserId(goodsItem.UserId)
		if err != nil {
			return nil, err
		}
		goods[i].Author.IsFollow, err = QueryFollowOrNot(currentId, userId)
		if err != nil {
			return nil, err
		}
		goods[i].PictureUrl = goodsItem.PictureUrl
		goods[i].Description = goodsItem.Description
		goods[i].Title = goodsItem.Title
		goods[i].IsFavorite = true
		//map中没有数据则自动为0
		goods[i].FavoriteCount = likeCountListMap[goodsItem.GoodsId]
		goods[i].CommentCount = commentCountListMap[goodsItem.GoodsId]
	}

	return
}
