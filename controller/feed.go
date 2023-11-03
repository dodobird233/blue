package controller

import (
	"blue/entity"
	"blue/global"
	"blue/service"
	"blue/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type feedResponse struct {
	Response  entity.Response
	GoodsList []entity.GoodsResponse `json:"goods_list,omitempty"`
	NextTime  int64                  `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	//获取 last_time, 找不到时使用当前时间
	LastTimeStr := c.DefaultQuery("last_time", "")
	var LastTime int64
	CurrentTime := time.Now().Unix()
	if LastTimeStr != "" {
		LastTimeTemp, err := time.Parse(time.RFC3339, LastTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status_code": -1, "status_msg": "Fail to get the last time."})
			return
		}
		LastTime = LastTimeTemp.Unix()
	} else {
		LastTime = CurrentTime
	}

	// 判断此时的商品列表是否为空
	var goodsList []entity.Goods
	var goodsIdList []int64
	var authorIdList []int64
	numGoods, err := service.GetNumGoods(&goodsList, &goodsIdList, &authorIdList, LastTime, global.MaxNumGoods)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": -1, "status_msg": "Fail to get the number of goods."})
		return
	}
	//fmt.Println(numGoods)

	// 如果是空的
	if numGoods == 0 {
		// 没有满足条件的商品
		c.JSON(http.StatusOK, feedResponse{
			Response:  entity.Response{StatusCode: 0, StatusMsg: "null"},
			GoodsList: nil,
			NextTime:  CurrentTime,
		})
		return
	}

	// 点赞信息获得
	LikeGoodsList, errLike := service.QueryLikeCountListByGoodsIdList(&goodsIdList)
	if errLike != nil {
		c.JSON(http.StatusNotFound, feedResponse{
			Response:  entity.Response{StatusCode: 1, StatusMsg: "Fail to get liked count for goods."},
			GoodsList: nil,
			NextTime:  LastTime,
		})
		return
	}
	//fmt.Println(LikeGoodsList)

	// 评论信息获得
	CommentGoodsList, errComment := service.QueryCommentCountListByGoodsIdList(&goodsIdList)
	if errComment != nil {
		c.JSON(http.StatusNotFound, feedResponse{
			Response:  entity.Response{StatusCode: 1, StatusMsg: "Fail to get comment count for goods."},
			GoodsList: nil,
			NextTime:  LastTime,
		})
		return
	}
	//fmt.Println(CommentGoodsList)

	// 点赞与否
	// 登录状态判断
	var userid int64
	isLogged := false
	token := c.PostForm("token")
	if token == "" {
		token = c.Query("token")
	}
	if token != "" {
		claims, errToken := util.Gettoken(token)
		if errToken == nil {
			userid, err = strconv.ParseInt(claims.UserId, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, feedResponse{
					Response:  entity.Response{StatusCode: 1, StatusMsg: err.Error()},
					GoodsList: nil,
					NextTime:  LastTime,
				})
				return
			}

			isLogged = true
		}
	}

	// 点赞与关注判断
	var isFavoriteList []bool
	var isFollowList []bool

	if isLogged {
		// 点赞列表
		isFavoriteList, err = service.ParseLikeGoodsListByUserIdFormGoodsId(userid, &goodsIdList)
		if err != nil {
			c.JSON(http.StatusBadRequest, feedResponse{
				Response:  entity.Response{StatusCode: 1, StatusMsg: err.Error()},
				GoodsList: nil,
				NextTime:  LastTime,
			})
			return
		}
		// 关注列表
		isFollowList, err = service.ParseFollowListByUserIdFormUserId(userid, &authorIdList)
		if err != nil {
			c.JSON(http.StatusBadRequest, feedResponse{
				Response:  entity.Response{StatusCode: 1, StatusMsg: err.Error()},
				GoodsList: nil,
				NextTime:  LastTime,
			})
			return
		}
	}

	isFavorite := false

	// 初始化列表信息
	var (
		goodsJsonList []entity.GoodsResponse
		goodsJson     entity.GoodsResponse
		author        entity.UserData
	)

	// 填充输出信息
	for i, goods := range goodsList {
		// author 获取
		author, err = service.UserInfoByUserId(authorIdList[i])
		if err != nil {
			//fmt.Println("Not found user")
			c.JSON(http.StatusNotFound, feedResponse{
				Response:  entity.Response{StatusCode: 1, StatusMsg: err.Error()},
				GoodsList: nil,
				NextTime:  LastTime,
			})
			return
		}

		// 登录时信息获取
		if isLogged {
			author.IsFollow = isFollowList[i]
			isFavorite = isFavoriteList[i]
		}

		// 信息填充
		goodsJson.Id = goods.GoodsId
		goodsJson.Author = author
		goodsJson.PictureUrl = goods.PictureUrl
		goodsJson.FavoriteCount = LikeGoodsList[i].LikeCnt
		goodsJson.CommentCount = CommentGoodsList[i].CommentCnt
		goodsJson.IsFavorite = isFavorite
		goodsJson.Title = goods.Title

		goodsJsonList = append(goodsJsonList, goodsJson)
	}

	nextTime := goodsList[numGoods-1].CreatedAt.Unix()
	// 输出商品流
	c.JSON(http.StatusOK, feedResponse{
		Response:  entity.Response{StatusCode: 0, StatusMsg: "null"},
		GoodsList: goodsJsonList,
		NextTime:  nextTime,
	})
}
