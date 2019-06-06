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
	e.POST("/user/newpwd",NewPwd)
	e.POST("/user/authkey",AuthKey)

	e.POST("/notif/new",NewNotif)
	e.GET("/notif/all",AllNotif)
	e.POST("/notif/remove",RemoveNotif)

	e.GET("/device/all",AllDevice)
	e.POST("/device/add",AddDevice)
	e.POST("/device/update",UpdateDevice)
	e.POST("/device/remove",RemoveDevice)
	e.POST("/device/code",FindDeviceCode)
	e.POST("/device/name",FindDeviceName)

	e.GET("/template/all",AllTem)
	e.POST("/template/add",AddTem)
	e.POST("/template/update",UpdateTem)
	e.POST("/template/remove",RemoveTem)
	e.POST("/template/key",FindTem)

	e.GET("/alarm/all",AllAlarm)

	e.GET("/alarmStg/all",AllStg)
	e.POST("/alarmStg/add",NewStg)
	e.POST("/alarmStg/update",UpdateStg)
	e.POST("/alarmStg/remove",RemoveStg)
	e.POST("/alarmStg/key",FindStg)

	return e
}