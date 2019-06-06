package deal

import (
	"fmt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewNotif(c *gin.Context){
	var notif model.Notif
	notif.Title="通知标题"
	notif.Content="通知内容"
	notif.Level=1
	notif.Type="1"
	err,_:=notif.Add()
	if err != nil {
		fmt.Println(err)
	}
}

func AllNotif(c *gin.Context){
	var  notif  model.Notif
	err,notifList:=notif.AllNotif()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":false,
			"msg":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"status":true,
		"data":notifList,
	})
}

func RemoveNotif(c *gin.Context){
	var  notif  model.Notif
	data:= struct {
		Key string `json:"key" form:"key"`
	}{}
	if err:=c.BindJSON(&data);err!=nil{
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":false,
			"msg":err.Error(),
		})
		return
	}
	notif.Key=data.Key
	if err:=notif.Remove();err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":false,
			"msg":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"status":true,
		"msg":"已删除通知",
	})
}
