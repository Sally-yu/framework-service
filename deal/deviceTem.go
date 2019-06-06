package deal

import (
	"fmt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AllTem(c *gin.Context) {
	tem := model.DeviceTemplate{}
	err, list := tem.All()
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

func RemoveTem(c *gin.Context) {
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
	tem := model.DeviceTemplate{}
	tem.Key = data.Key
	if err := tem.Remove(); err != nil {
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

func UpdateTem(c *gin.Context) {
	tem := model.DeviceTemplate{}
	if err := c.Bind(&tem); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	if err := tem.Update(); err != nil {
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

func AddTem(c *gin.Context) {
	tem := model.DeviceTemplate{}
	if err := c.Bind(&tem); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	if err := tem.FindByCode(); err == nil {
		fmt.Println("同编号设备已存在！")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "同编号设备已存在",
		})
		return
	}
	if err := tem.New(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "已添加设备模板",
		})
	}
}

func FindTem(c *gin.Context) {
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
	tem := model.DeviceTemplate{}
	tem.Key = data.Key
	if err := tem.Find(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   tem,
		})
	}
}
