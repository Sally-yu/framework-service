package deal

import (
	"fmt"
	"framework-service/jwt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

//请求所有的数据库信息
func AlldbMgr(context *gin.Context) {
	dbMgr := model.DbMgr{}
	err, list := dbMgr.All()
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   list,
		})
	}
}
//添加一条记录
func AddDbMgr(context *gin.Context) {
	claims := context.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
		dbMgr := model.DbMgr{}
		if err := context.Bind(&dbMgr); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
			b, msg := dbMgr.Insert()
		if !b {
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    msg,
			})
			return
		}
			context.JSON(http.StatusOK, gin.H{
				"status": true,
				"msg":    msg,
			})

	}
}
//更新一条记录
func UpdateDbMgr(context *gin.Context) {
	claims := context.MustGet("claims").(*jwt.CustomClaims) //需要携带token
	if claims != nil {
		dbMgr := model.DbMgr{}
		if err := context.BindJSON(&dbMgr); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		if err := dbMgr.Update(); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{
				"status": true,
				"msg":    "数据库信息已更新",
			})
		}
	}
}
//删除一条记录
func DeleteDbMgr(context *gin.Context) {
	claims := context.MustGet("claims").(*jwt.CustomClaims) //需要携带token
	if claims != nil {
		data := struct {
			Serverip string `json:"serverip" form:"serverip"`
		}{}
		if err := context.Bind(&data); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		dbMgr := model.DbMgr{}
		dbMgr.Serverip = data.Serverip
		if err := dbMgr.Delete(); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{
				"status": true,
				"msg":    "已删除",
			})
		}
	}
}

//查找数据库
func FindServerIp(context *gin.Context) {
	claims := context.MustGet("claims").(*jwt.CustomClaims) //需要携带token
	if claims != nil {
		data := struct {
			Serverip string `json:"serverip" form:"serverip"`
		}{}
		d := model.DbMgr{}
		if err := context.Bind(&data); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		d.Serverip = data.Serverip
		if err := d.FindByServerip(); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{
				"status": true,
				"data":   d,
			})
		}
	}
}

//测试数据库的连接
func TestPing(context *gin.Context) {
	claims := context.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
		dbMgr := model.DbMgr{}
		if err := context.Bind(&dbMgr); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		if err := dbMgr.TestPing(); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    "测试连接失败！",
			})
			//return
		} else {
			context.JSON(http.StatusOK, gin.H{
				"status": true,
				"msg":   "测试连接成功!",
			})
		}
	}
}