package main

import (
	"fmt"
	. "framework-service/deal" //合并导入，不需要包名引用
	"framework-service/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)


func main() {
	engine:=gin.Default()
	//config:=cors.Default() //源码指示允许cors
	//engine.Use(cors.Default()) //允许cors
	engine.Use(Cors())
	engine=Handle(engine)
	engine.Run(":9060")
}

func Handle(e *gin.Engine) *gin.Engine{
	e.GET("/get",Get)
	e.GET("/shell",Shell)

	e.GET("/rsakey",RsaKey)
	e.POST("/user/login",Login)

	user:=e.Group("/user")
	user.Use(jwt.JWTAuth())
	{
		user.POST("/user/key", FindUser)
		user.GET("/user/all", AllUser)
		user.POST("/user/add", AddUser)
		user.POST("/user/remove", RemoveUser)
		user.POST("/user/update", UpdateUser)
		user.POST("/user/newpwd", NewPwd)
		user.POST("/user/authkey", AuthKey)
	}

	e.POST("/notif/new",NewNotif)
	e.GET("/notif/all",AllNotif)
	e.POST("/notif/remove",RemoveNotif)
	e.POST("/notif/update",UpdateNotif)

	devcie := e.Group("/device") //需要验证token
	devcie.Use(jwt.JWTAuth())
	{
		devcie.GET("/all", AllDevice)
		devcie.POST("/add",AddDevice)
		devcie.POST("/update",UpdateDevice)
		devcie.POST("/remove",RemoveDevice)
	}

	e.POST("/device/code",FindDeviceCode)
	e.POST("/device/name",FindDeviceName)
	e.POST("/device/value",DeviceValue)

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


func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method      //请求方法
		origin := c.Request.Header.Get("Origin")        //请求头部
		var headerKeys []string                             // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")        // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")      //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")      // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")        // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")       //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")       // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()        //  处理请求
	}
}