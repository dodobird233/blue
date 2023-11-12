package service

import (
	"blue/entity"
	"blue/global"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 查询商品id列表
func QueryGoodsIdListByUserId(userId int64) (goodsIdList []int64, err error) {
	result := global.DB.Model(&entity.Goods{}).Select("goods_id").Where("user_id = ?", userId).Find(&goodsIdList)
	if result.Error != nil {
		err = result.Error
		return nil, err
	}
	return
}

// 查询商品对象列表
func QueryGoodsListByUserId(userId int64) (goodsList []entity.Goods, err error) {
	if global.DB.Where("user_id = ?", userId).Find(&goodsList).Error != nil {
		return
	}
	return
}

// 查询商品封装返回对象列表
func GetPostGoodsListByUserId(userId int64) (goods []entity.GoodsResponse, err error) {
	//查询商品对象列表
	goodsList, err := QueryGoodsListByUserId(userId)
	if err != nil {
		return nil, err
	}
	//构造商品id列表
	goodsIdList := make([]int64, len(goodsList))
	for i, goodsItem := range goodsList {
		goodsIdList[i] = goodsItem.GoodsId
	}
	//根据商品id列表查询点赞数量
	likeCountList, err := QueryLikeCountListByGoodsIdList(&goodsIdList)
	if err != nil {
		return nil, err
	}
	likeCountListMap := map[int64]int64{}
	for _, likeCount := range likeCountList {
		likeCountListMap[likeCount.GoodsId] = likeCount.LikeCnt
	}
	//根据商品id列表查询评论数量
	commentCountList, err := QueryCommentCountListByGoodsIdList(&goodsIdList)
	if err != nil {
		return nil, err
	}
	commentCountListMap := map[int64]int64{}
	for _, likeCount := range commentCountList {
		commentCountListMap[likeCount.GoodsId] = likeCount.CommentCnt
	}
	goods = make([]entity.GoodsResponse, len(goodsList))
	for i, goodsItem := range goodsList {
		goods[i].Id = goodsItem.GoodsId
		goods[i].Author, err = UserInfoByUserId(goodsItem.UserId)
		if err != nil {
			return nil, err
		}
		//仅有登录用户自己
		goods[i].Author.IsFollow, err = QueryFollowOrNot(userId, userId)
		if err != nil {
			return nil, err
		}
		goods[i].PictureUrl = goodsItem.PictureUrl
		goods[i].Description = goodsItem.Description
		goods[i].Title = goodsItem.Title
		goods[i].IsFavorite = true
		goods[i].FavoriteCount = likeCountListMap[goodsItem.GoodsId]
		goods[i].CommentCount = commentCountListMap[goodsItem.GoodsId]
	}
	return
}

func SavePictureUrls(savePath []string, pictureName []string) (string, error) {
	// 存储空间名称
	bucketName := "blue-project"
	urls := ""
	for i := range savePath {
		// Object完整路径，完整路径不包含Bucket名称
		objectName := pictureName[i]
		// 本地文件的完整路径
		localFileName := savePath[i]

		// 创建OSSClient实例
		client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com", global.ACCESS_KEY_ID, global.ACCESS_SEC)
		if err != nil {
			return "", err
		}

		// 获取存储空间
		bucket, err := client.Bucket(bucketName)
		if err != nil {
			return "", err
		}
		// 上传文件
		err = bucket.PutObjectFromFile(objectName, localFileName)
		if err != nil {
			return "", err
		}
		urls += "https://blue-project.oss-cn-hangzhou.aliyuncs.com/" + pictureName[i]
		if i != len(savePath)-1 {
			urls += ","
		}
	}
	return urls, nil
}
