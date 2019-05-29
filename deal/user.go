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
	Data string `json:"data" form:"data"`
}

//登录验证
func Login(c *gin.Context) {
	user,err:=Decrypt(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	b, msg := user.Auth()
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
	})

}

//注册
func AddUser(c *gin.Context) {
	user,err:=Decrypt(c)
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

//更新用户信息
func UpdateUser(c *gin.Context) {
	user,err:=Decrypt(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
	}
}

//移除用户
func RemoveUser(c *gin.Context) {
	user,err:=Decrypt(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	if err := user.Remove(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": false,
		"msg":    "已移除用户",
	})
}

func Decrypt(c *gin.Context) (model.User,error)  {
	var user model.User
	data := EncryptData{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return user,err
	}
	destr, err := gorsa.PriKeyDecrypt(data.Data, Pirvatekey) //解析出[]byte 可解析出原json对象
	if err != nil {
		fmt.Println(err.Error())
		return user,err
	}
	err = json.Unmarshal(destr, &user) //解析json
	if err != nil {
		fmt.Println(err.Error())
		return user,err
	}
	return user,nil
}