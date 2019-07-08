package deal

import (
	"fmt"
	"framework-service/jwt"
	"framework-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpsertModel(c *gin.Context) {
		list := []model.Model3d{}
		if err := c.BindJSON(&list); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		for i := 0; i < len(list); i++ {
			list[i].Insert()
		}
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   list,
		})
}

func AllModel(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //header携带token
	if claims != nil {
		m := model.Model3d{}
		err, list := m.FindAll()
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
}

func ModelRel(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims) //需要携带token
	if claims != nil {
		m := model.Model3d{}
		if err := c.Bind(&m); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		}
		fmt.Println(m)
		if err := m.Update(); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"msg":    err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"msg":    "信息已更新",
			})
		}
	}
}
