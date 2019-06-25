package deal

import (
	"fmt"
	"framework-service/jwt"
	"framework-service/model"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// LoginResult 登录结果结构
type LoginResult struct {
	Token string     `json:"token" form:"token"`
	User  model.User `json:"user" form:"token"`
}

// 生成令牌
func generateToken(c *gin.Context, user model.User) {
	j := &jwt.JWT{
		[]byte(jwt.GetSignKey()),
	}
	claims := jwt.CustomClaims{
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),    // 签名生效时间
			//ExpiresAt: int64(time.Now().Unix() + 10000),
			ExpiresAt: int64(time.Now().Unix() + 2592000), // 过期时间 30天
			Issuer:    "Inspur-MOM",                       //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    "生成token失败",
		})
		return
	}

	u:=model.UserToken{}
	u.User=user.Key
	u.Token=token
	u.Update()
	user.Pwd=""
	data := LoginResult{
		Token: token,
		User:  user,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "用户验证成功！",
		"data":   data,
	})
	return
}

func tokenValidate(c *gin.Context) {
	data := struct {
		Token string `json:"token" form:"token"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    err.Error()})
		return
	}
	token := data.Token
	j := &jwt.JWT{
		[]byte(jwt.GetSignKey()),
	}
	_, err := j.ParseToken(token)
	switch err {
	case nil:
		fmt.Println("token验证通过")
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "token验证通过",
		})
		break
	case jwt.TokenExpired:
		fmt.Println("token已过期")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "token已过期",
		})
		break
	case jwt.TokenInvalid:
		fmt.Println("token验证失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "token验证失败",
		})
		break
	case jwt.TokenNotValidYet:
		fmt.Println("token已过期")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "token已过期",
		})
		break
	case jwt.TokenMalformed:
		fmt.Println("token验证失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "token验证失败",
		})
		break
	default:
		break
	}
}
