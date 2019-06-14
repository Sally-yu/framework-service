package deal

import (
	"fmt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

func AllDevice(c *gin.Context) {
	device := model.Device{}
	err, list := device.All()
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

func RemoveDevice(c *gin.Context) {
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
	device := model.Device{}
	device.Key = data.Key
	if err := device.Remove(); err != nil {
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

func UpdateDevice(c *gin.Context) {
	device := model.Device{}
	if err := c.BindJSON(&device); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	if err := device.Update(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "设备信息已更新",
		})
	}
}

func AddDevice(c *gin.Context) {
	device := model.Device{}
	if err := c.Bind(&device); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	if err := device.FindByCode(); err == nil {
		fmt.Println("设备编号已被使用！")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "设备编号已被使用",
		})
		return
	}
	if err := device.FindByName(); err == nil {
		fmt.Println("同名设备已存在！")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "同名设备已存在",
		})
		return
	}
	if b, msg := device.Insert(); !b {
		fmt.Println(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    msg,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    msg,
		})
	}
}

func FindDeviceCode(c *gin.Context) {
	data := struct {
		Code string `json:"code" form:"code"`
	}{}
	d := model.Device{}
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	d.Code = data.Code
	if err := d.FindByCode(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   d,
		})
	}
}

func FindDeviceName(c *gin.Context) {
	data := struct {
		Name string `json:"name" form:"name"`
	}{}
	d := model.Device{}
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	d.Name = data.Name
	if err := d.FindByName(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   d,
		})
	}
}

type Res struct {
	Device string     `json:"device" form:"device"`
	Data   [] ResData `json:"data" form:"data"`
}

type ResData struct {
	AttCode string      `json:"attcode" form:"attcode"`
	Value   interface{} `json:"value" form:"value"`
}

//设备属性值
func DeviceValue(c *gin.Context) {
	data := struct {
		Keys []string `json:"keys" form:"keys"`
	}{}
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	res := []Res{}
	for k := range data.Keys {
		r := GetAttValue(data.Keys[k])
		res = append(res, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   res,
	})
}

func GetAttValue(key string) Res {
	r := Res{}
	r.Device = key
	device := model.Device{}
	device.Key = key
	device.Find()
	for i := range device.Attrs {
		data:=ResData{}
		data.AttCode=device.Attrs[i].Code
		data.Value=rand.Intn(100)
		r.Data = append(r.Data,data)
	}
	return r
}
