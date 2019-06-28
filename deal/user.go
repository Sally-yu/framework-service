package deal

import (
	"encoding/json"
	"fmt"
	"framework-service/crypt"
	"framework-service/jwt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"github.com/wenzhenxi/gorsa"
	"net/http"
)

type EncryptData struct {
	User string `json:"user" form:"user"`
}

//登录验证，刷新token
func Login(c *gin.Context) {
	data := struct {
		Name string `json:"name" form:"name"`
		Pwd  string `json:"pwd" form:"pwd"` //加密过的密码 密文传输
	}{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error()})
		return
	}
	user := model.User{}
	user.Uname = data.Name
	if err := user.FindByName(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"msg":    "密码错误或用户信息不存在"})
		return
	}
	if !user.ComparePwd(data.Pwd) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"msg":    "密码错误或用户信息不存在"})
		return
	}
	generateToken(c, user) //生成token
	//c.JSON(http.StatusOK, gin.H{
	//	"status": true,
	//	"msg":    "验证成功",
	//	"data":   user.Key,
	//})
}

//校验用户key和密码是否对应
func AuthKey(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
		data := struct {
			Key string `json:"key" form:"key"`
			Pwd string `json:"pwd" form:"pwd"`
		}{}
		if err := c.BindJSON(&data); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error()})
			return
		}
		user := model.User{}
		user.Key = data.Key
		if err := user.FindUser(); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    err.Error()})
			return
		}
		if !user.ComparePwd(data.Pwd) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    "密码错误或用户信息不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "验证成功",
			"data":   user.Key,
		})
	}
}

//注册
func AddUser(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {

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
}

//查找用户
func FindUser(c *gin.Context) {
	var user model.User
	var data struct {
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
	user.Key = data.Key
	err := user.FindUser()
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
		"status": true,
		"data":   user,
	})
}

//查找用户
func FindPhone(c *gin.Context) {
	var user model.User
	var data struct {
		Phone string `json:"phone" form:"phone"`
	}
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	fmt.Println(data.Phone)
	user.Phone = data.Phone
	err := user.FindByPhone()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "未能找到该手机号码的注册信息",
		})
		return
	}
	user.Pwd = ""
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}

//get所有用户
func AllUser(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
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
}

//更新用户信息
func UpdateUser(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
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
}

//更新用户密码
func NewPwd(c *gin.Context) {
	data := struct {
		Key string `json:"key" form:"key"`
		Pwd string `json:"pwd" form:"pwd"` //加密过的密码 直接传直接存 不解密
	}{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error()})
		return
	}
	user := model.User{}
	user.Key = data.Key
	if err := user.FindUser(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error()})
		return
	}
	user.Pwd = data.Pwd
	if err := user.Update(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "密码已更新"})
	return
}

//移除用户
func RemoveUser(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
		var key struct {
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
		user.Key = key.Key
		if err := user.Remove(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "已删除用户",
		})
	}
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
	destr, err := gorsa.PriKeyDecrypt(data.User, crypt.Pirvatekey) //解析出[]byte 可解析出原json对象
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	fmt.Println(destr)
	err = json.Unmarshal([]byte(destr), &user) //解析json
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	return user, nil
}

func RsaKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   crypt.PublicKey,
	})
}
