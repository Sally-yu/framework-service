package deal

import (
"fmt"
"framework-service/model"
"github.com/gin-gonic/gin"
"net/http"
)

func AllAlarm(c *gin.Context) {
	alarm := model.Alarm{}
	err, list := alarm.All()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   list,
		})
	}
}