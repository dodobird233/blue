package controller

import (
	"blue/entity"
	"blue/global"
	"blue/service"
	"blue/util"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	// 获取商品唯一标识 id
	node, err := snowflake.NewNode(1)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{"status_code": 1, "status_msg": "failed to generate snowflake"})
		c.JSON(http.StatusBadRequest, entity.Response{
			StatusCode: 1,
			StatusMsg:  "failed to generate snowflake for goods",
		})
	}
	var pictureName []string
	var savePath []string

	goodsId := node.Generate().Int64()              //生成唯一id
	name := strconv.FormatUint(uint64(goodsId), 10) //生成唯一name

	// todo begin
	// 获取文件
	r := c.Request
	//设置内存大小
	r.ParseMultipartForm(32 << 20)
	//获取上传的文件组
	files := r.MultipartForm.File["picture"]
	for i := 0; i < len(files); i++ {
		//打开上传文件
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		pictureName = append(pictureName, name+files[i].Filename)
		savePath = append(pictureName, filepath.Join("./public/goods/", name+files[i].Filename))
		err = c.SaveUploadedFile(files[i], filepath.Join("./public/goods/", name+files[i].Filename)) //库函数保存文件到public目录下
		if err != nil {
			c.JSON(http.StatusInternalServerError, entity.Response{
				StatusCode: -1,
				StatusMsg:  "fail to save the file to the path.",
			})
			return
		}
		//debug
		fmt.Println(files[i].Filename) //输出上传的文件名
	}
	// todo end
	// 商品图片存入oss,返回拼接的url
	urls, err := service.SavePictureUrls(savePath, pictureName) // ok
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{
			StatusCode: -1,
			StatusMsg:  "fail to upload pi",
		})
		return
	}
	for _, pa := range savePath {
		err = os.Remove(pa)
		if err != nil {
			c.JSON(http.StatusInternalServerError, entity.Response{
				StatusCode: -1,
				StatusMsg:  "fail to delete the file.",
			})
			return
		}
	}

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
