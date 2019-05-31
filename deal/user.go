package deal

import (
	"encoding/json"
	"fmt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"github.com/wenzhenxi/gorsa"
	"net/http"
)

type EncryptData struct {
	User string `json:"user" form:"user"`
}

//登录验证
func Login(c *gin.Context) {
	user, err := Decrypt(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	b, msg, key := user.Auth()
	if !b { //验证不通过
		c.JSON(http.StatusBadRequest, gin.H{
			"status": b,
			"msg":    msg,
		})
		return
	}
	err = user.Login()
	if err != nil { //更新登录时间
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": b,
		"msg":    "验证成功！",
		"data":   key,
	})

}

//注册
func AddUser(c *gin.Context) {
	user, err := Decrypt(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	b, msg := user.Insert()
	fmt.Println(msg)
	if b {
		c.JSON(http.StatusOK, gin.H{
			"status": b,
			"msg":    msg,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": b,
			"msg":    msg,
		})
	}
}

//查找用户
func FindUser(c *gin.Context) {
	var user model.User
	var data struct{
		Key string `json:"key" form:"key"`
	}
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	fmt.Println(data.Key)
	user.Key=data.Key
	err:=user.FindUser()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	user.Pwd = ""
	c.JSON(http.StatusOK, gin.H{
		"status":true,
		"data":user,
	})
}

//get所有用户
func AllUser(c *gin.Context) {
	var user model.User
	err, res := user.All()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    nil,
			"data":   res,
		})
	}
}

//更新用户信息
func UpdateUser(c *gin.Context) {
	user, err := Decrypt(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "已更新用户信息",
	})
}

//移除用户
func RemoveUser(c *gin.Context) {
	var key struct{
		Key string `json:"key" form:"key"`
	}
	if err := c.BindJSON(&key); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error()})
		return
	}
	var user model.User
	user.Key=key.Key
	if err := user.Remove(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "已移除用户",
	})
}

//解密用户信息
func Decrypt(c *gin.Context) (model.User, error) {
	var user model.User
	data := EncryptData{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error()})
		return user, err
	}
	destr, err := gorsa.PriKeyDecrypt(data.User, Pirvatekey) //解析出[]byte 可解析出原json对象
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	fmt.Println(destr)
	err = json.Unmarshal(destr, &user) //解析json
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	return user, nil
}
