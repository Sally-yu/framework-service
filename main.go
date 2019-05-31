package main

import (
	. "framework-service/deal" //合并导入，不需要包名引用
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func main() {
	engine:=gin.Default()
	//config:=cors.Default() //源码指示允许cors
	engine.Use(cors.Default()) //允许cors
	engine=Handle(engine)
	engine.Run(":9060")
}

func Handle(e *gin.Engine) *gin.Engine{
	e.GET("/get",Get)
	e.GET("/shell",Shell)

	e.POST("/user/key",FindUser)
	e.GET("/user/all",AllUser)
	e.POST("/user/login",Login)
	e.POST("/user/add",AddUser)
	e.POST("/user/remove",RemoveUser)
	e.POST("/user/update",UpdateUser)

	return e
}