package service

import "blue/entity"

// 在一个给定的 int64 数组中查找给定元素
func FindInt64(target int64, intArr *[]int64) bool {
	for _, element := range *intArr {
		if target == element {
			return true
		}
	}
	return false
}

// 在一个给定 GoodsLikeCnt 列表中查找给定商品 id 是否存在，不存在返回 0，存在返回点赞值
func FindGoodsIdFromGoodsLikeCntList(goodsId int64, likeCountList *[]entity.GoodsLikeCnt) int64 {
	for _, element := range *likeCountList {
		if goodsId == element.GoodsId {
			return element.LikeCnt
		}
	}
	return 0
}

// 在一个给定 GoodsCommentCnt 列表中查找给定商品 id 是否存在，不存在返回 0，存在返回评论值
func FindGoodsIdFromGoodsCommentCntList(goodsId int64, commentCountList *[]entity.GoodsCommentCnt) int64 {
	for _, element := range *commentCountList {
		if goodsId == element.GoodsId {
			return element.CommentCnt
		}
	}
	return 0
}
