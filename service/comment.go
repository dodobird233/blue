package service

import (
	"blue/entity"
	"blue/global"
	"blue/util"
	"errors"
	"time"
)

// 根据商品id列表查询评论数量列表
func QueryCommentCountListByGoodsIdList(goodsIdList *[]int64) ([]entity.GoodsCommentCnt, error) {
	var getCommentCountList []entity.GoodsCommentCnt
	result := global.DB.Model(&entity.Comment{}).Select("goods_id", "count(goods_id) as comment_cnt").Where("goods_id in ?", *goodsIdList).Group("goods_id").Find(&getCommentCountList)
	if result.Error != nil {
		err := errors.New("commentList query failed")
		return nil, err
	}
	// 找数据找齐了
	if len(*goodsIdList) == len(getCommentCountList) {
		return getCommentCountList, nil
	}
	// 数据不全，误差部分补全为0
	var commentCountList []entity.GoodsCommentCnt
	commentCountList = make([]entity.GoodsCommentCnt, len(*goodsIdList))
	for i, goodsId := range *goodsIdList {
		commentCountList[i].GoodsId = goodsId
		commentCountList[i].CommentCnt = FindGoodsIdFromGoodsCommentCntList(goodsId, &getCommentCountList)
	}
	return commentCountList, nil
}

// 增加评论
func AddComment(currentId int64, goodsId int64, commentText string) (err error) {
	var addComment entity.Comment
	addComment.CommentId = util.GetNextId()
	addComment.UserId = currentId
	addComment.GoodsId = goodsId
	addComment.Content = commentText
	addComment.CreateDate = time.Now().Format("01-02")
	result := global.DB.Model(&entity.Comment{}).Create(&addComment)
	if result.Error != nil {
		return err
	}
	return
}

// 删除评论
func CancelComment(currentId int64, goodsId int64, commentId int64) (err error) {
	var cancelComment entity.Comment
	cancelComment.CommentId = commentId
	cancelComment.UserId = currentId
	cancelComment.GoodsId = goodsId
	result := global.DB.Model(&entity.Comment{}).Delete(&cancelComment)
	if result.Error != nil {
		return err
	}
	if result.RowsAffected == 0 {
		err = errors.New("comment not found")
		return err
	}
	return
}

func GetCommentListByGoodsId(currentId int64, goodsId int64) (comments []entity.CommentResponse, err error) {
	var commentList []entity.Comment
	if global.DB.Model(&entity.Comment{}).Where("goods_id=?", goodsId).Find(&commentList).Error != nil {
		return
	}
	comments = make([]entity.CommentResponse, len(commentList))
	for i, comment := range commentList {
		comments[i].Id = comment.CommentId
		comments[i].User, err = UserInfoByUserId(comment.UserId)
		if err != nil {
			return nil, err
		}
		comments[i].User.IsFollow, err = QueryFollowOrNot(currentId, comment.UserId)
		if err != nil {
			return nil, err
		}
		comments[i].Content = comment.Content
		comments[i].CreateDate = comment.CreateDate
	}
	return
}
