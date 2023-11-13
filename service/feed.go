package service

import (
	"blue/entity"
	"blue/global"
)

// GetNumGoods 获取商品列表中符合结果的商品数据，以及相应的商品，作者信息列表
func GetNumGoods(goods *[]entity.Goods, goodsIdList *[]int64, AuthorIdList *[]int64, LastTime int64, MaxNumGoods int) (int, error) {
	goodsList, _ := GetGoodsListFromRedis()
	if goodsList != nil {
		numGoods := len(*goodsList)
		*goodsIdList = make([]int64, numGoods)
		*AuthorIdList = make([]int64, numGoods)
		*goods = make([]entity.Goods, numGoods)
		for i, goodsItem := range *goodsList {
			(*AuthorIdList)[i] = goodsItem.UserId
			(*goodsIdList)[i] = goodsItem.GoodsId
			(*goods)[i] = goodsItem
		}
		return numGoods, nil
	}
	//查询数据库
	global.DB.Order("created_at desc").Limit(MaxNumGoods).Find(&goods)
	if goods == nil {
		return 0, nil
	}
	numGoods := len(*goods)
	// 统计作者 id 以及商品 id
	*goodsIdList = make([]int64, numGoods)
	*AuthorIdList = make([]int64, numGoods)
	for i, goodsItem := range *goods {
		(*AuthorIdList)[i] = goodsItem.UserId
		(*goodsIdList)[i] = goodsItem.GoodsId
	}
	//加入缓存
	AddGoodsByGoodsIdFromCacheToRedis(goods)
	//do not handle
	return numGoods, nil
}
