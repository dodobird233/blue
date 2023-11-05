package controller

import (
	"blue/entity"
	"blue/global"
	"blue/service"
	"blue/util"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type goodsListResponse struct {
	entity.Response
	GoodsList []entity.GoodsResponse `json:"goods_list"`
}

func Publish(c *gin.Context) {
	//校验token并获取当前用户id
	token := c.PostForm("token")
	claims, _ := util.Gettoken(token)
	userid, _ := strconv.ParseInt(claims.UserId, 10, 64)
	title := c.PostForm("title")
	desc := c.PostForm("description")
	urls := c.PostForm("picture_urls")
	// 获取文件

	// 获取商品唯一标识 id
	node, err := snowflake.NewNode(1)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{"status_code": 1, "status_msg": "failed to generate snowflake"})
		c.JSON(http.StatusBadRequest, entity.Response{
			StatusCode: 1,
			StatusMsg:  "failed to generate snowflake for goods",
		})
	}
	goodsId := node.Generate().Int64()

	// 商品图片存入数据库
	service.SavePictureUrls(urls)

	goods := entity.Goods{
		GoodsId:     goodsId,
		PictureUrl:  urls,
		Description: desc,
		Title:       title,
		UserId:      userid,
	}
	err = global.DB.Create(&goods).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{
			StatusCode: -1,
			StatusMsg:  "fail to add the goods into SQL",
		})
		return
	}

	//fmt.Printf("写入数据库\n")

	c.JSON(http.StatusOK, entity.Response{
		StatusCode: 0,
		StatusMsg:  "null",
	})
}

func PublishList(c *gin.Context) {
	//校验token并获取当前用户id
	token := c.Query("token")
	_, err := util.Gettoken(token)
	if err != nil {
		c.JSON(http.StatusOK, entity.Response{StatusCode: 1, StatusMsg: "token error"})
		return
	}
	uid, _ := strconv.Atoi(c.Query("user_id"))
	goods, err := service.GetPostGoodsListByUserId(int64(uid))
	if err != nil {
		c.JSON(http.StatusOK, entity.Response{StatusCode: 2, StatusMsg: "get liked goods list failed"})
		return
	}
	//封装返回
	c.JSON(http.StatusOK, goodsListResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		GoodsList: goods,
	})
}
