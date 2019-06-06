package deal

import (
	"fmt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewStg(c *gin.Context)  {
	alarmStg := model.AlarmStg{}
	if err := c.Bind(&alarmStg); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	if err := alarmStg.Find(); err == nil {
		fmt.Println("同编号策略已存在！")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "同编号策略已存在",
		})
		return
	}
	if err := alarmStg.Insert(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "已添加报警策略",
		})
	}
}

func AllStg(c *gin.Context){
	alarmStg := model.AlarmStg{}
	err, list := alarmStg.All()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   list,
		})
	}
}

func FindStg(c *gin.Context){
	data := struct {
		Key string `json:"key" form:"key"`
	}{}
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	alarmStg := model.AlarmStg{}
	alarmStg.Key = data.Key
	if err := alarmStg.Find(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   alarmStg,
		})
	}
}

func UpdateStg(c *gin.Context)  {
	alarmStg := model.AlarmStg{}
	if err := c.Bind(&alarmStg); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	if err := alarmStg.Update(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "已更新",
		})
	}
}

func RemoveStg(c *gin.Context)  {
	data := struct {
		Key string `json:"key" form:"key"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	alarmStg := model.AlarmStg{}
	alarmStg.Key = data.Key
	if err := alarmStg.Remove(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "已删除",
		})
	}
}
