package deal

import (
	myjwt "framework-service/jwt"
	"framework-service/model"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
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
	j := &myjwt.JWT{
		[]byte(myjwt.GetSignKey()),
	}
	claims := myjwt.CustomClaims{
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 2592000), // 过期时间 30天
			Issuer:    "Inspur-MOM",                    //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	log.Println(token)
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
