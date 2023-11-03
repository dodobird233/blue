package service

import (
	"blue/entity"
	"blue/global"
)

// GetNumGoods 获取商品列表中符合结果的商品数据，以及相应的商品，作者信息列表
func GetNumGoods(goods *[]entity.Goods, goodsIdList *[]int64, AuthorIdList *[]int64, LastTime int64, MaxNumGoods int) (int, error) {
	query := global.DB.Order("created_at desc").
		Limit(MaxNumGoods).
		Where("UNIX_TIMESTAMP(created_at) <= ?", LastTime)
	query.Find(goods)

	numGoods := len(*goods)

	// 统计作者 id 以及商品 id
	*AuthorIdList = make([]int64, numGoods)
	*goodsIdList = make([]int64, numGoods)
	for i, goodsItem := range *goods {
		(*AuthorIdList)[i] = goodsItem.UserId
		(*goodsIdList)[i] = goodsItem.GoodsId
	}

	return numGoods, nil
}
