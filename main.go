/**
 *            .,,       .,:;;iiiiiiiii;;:,,.     .,,
 *          rGB##HS,.;iirrrrriiiiiiiiiirrrrri;,s&##MAS,
 *         r5s;:r3AH5iiiii;;;;;;;;;;;;;;;;iiirXHGSsiih1,
 *            .;i;;s91;;;;;;::::::::::::;;;;iS5;;;ii:
 *          :rsriii;;r::::::::::::::::::::::;;,;;iiirsi,
 *       .,iri;;::::;;;;;;::,,,,,,,,,,,,,..,,;;;;;;;;iiri,,.
 *    ,9BM&,            .,:;;:,,,,,,,,,,,hXA8:            ..,,,.
 *   ,;&@@#r:;;;;;::::,,.   ,r,,,,,,,,,,iA@@@s,,:::;;;::,,.   .;.
 *    :ih1iii;;;;;::::;;;;;;;:,,,,,,,,,,;i55r;;;;;;;;;iiirrrr,..
 *   .ir;;iiiiiiiiii;;;;::::::,,,,,,,:::::,,:;;;iiiiiiiiiiiiri
 *   iriiiiiiiiiiiiiiii;;;::::::::::::::::;;;iiiiiiiiiiiiiiiir;
 *  ,riii;;;;;;;;;;;;;:::::::::::::::::::::::;;;;;;;;;;;;;;iiir.
 *  iri;;;::::,,,,,,,,,,:::::::::::::::::::::::::,::,,::::;;iir:
 * .rii;;::::,,,,,,,,,,,,:::::::::::::::::,,,,,,,,,,,,,::::;;iri
 * ,rii;;;::,,,,,,,,,,,,,:::::::::::,:::::,,,,,,,,,,,,,:::;;;iir.
 * ,rii;;i::,,,,,,,,,,,,,:::::::::::::::::,,,,,,,,,,,,,,::i;;iir.
 * ,rii;;r::,,,,,,,,,,,,,:,:::::,:,:::::::,,,,,,,,,,,,,::;r;;iir.
 * .rii;;rr,:,,,,,,,,,,,,,,:::::::::::::::,,,,,,,,,,,,,:,si;;iri
 *  ;rii;:1i,,,,,,,,,,,,,,,,,,:::::::::,,,,,,,,,,,,,,,:,ss:;iir:
 *  .rii;;;5r,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,sh:;;iri
 *   ;rii;:;51,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,.:hh:;;iir,
 *    irii;::hSr,.,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,.,sSs:;;iir:
 *     irii;;:iSSs:.,,,,,,,,,,,,,,,,,,,,,,,,,,,..:135;:;;iir:
 *      ;rii;;:,r535r:...,,,,,,,,,,,,,,,,,,..,;sS35i,;;iirr:
 *       :rrii;;:,;1S3Shs;:,............,:is533Ss:,;;;iiri,
 *        .;rrii;;;:,;rhS393S55hh11hh5S3393Shr:,:;;;iirr:
 *          .;rriii;;;::,:;is1h555555h1si;:,::;;;iirri:.
 *            .:irrrii;;;;;:::,,,,,,,,:::;;;;iiirri:.
 *               .:irrriiiiiii;;;;;;;;iiiiiiirri:,.
 *                  .,;;iirrrrrrrrrrrrrrrrri;:.
 *                         ..,:::;;;;:::,,.
**/

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
	go FileServer()//并发一个文件服务
	engine := gin.Default()
	//config:=cors.Default() //源码指示允许cors
	//engine.Use(cors.Default()) //允许cors
	engine.Use(Cors()) //需要header携带token，gin默认的跨域不能支持，手写
	engine = Handle(engine)
	engine.Run(":9060")
}

func Handle(e *gin.Engine) *gin.Engine {
	e.GET("/get", Get)
	e.GET("/shell", Shell)

	e.GET("/rsakey", RsaKey)
	e.POST("/login", Login)
	e.POST("/phone", FindPhone)
	e.POST("/newpwd", NewPwd)

	user := e.Group("/user")
	user.Use(jwt.JWTAuth())
	{
		user.POST("/key", FindUser)
		user.GET("/all", AllUser)
		user.POST("/add", AddUser)
		user.POST("/remove", RemoveUser)
		user.POST("/update", UpdateUser)
		user.POST("/authkey", AuthKey)
	}

	notif := e.Group("/notif")
	notif.Use(jwt.JWTAuth())
	{
		notif.POST("/new", NewNotif)
		notif.GET("/all", AllNotif)
		notif.POST("/remove", RemoveNotif)
		notif.POST("/update", UpdateNotif)
	}
	device := e.Group("/device") //需要验证token
	device.Use(jwt.JWTAuth())
	{
		device.GET("/all", AllDevice)
		device.POST("/add", AddDevice)
		device.POST("/update", UpdateDevice)
		device.POST("/remove", RemoveDevice)
		device.POST("/code", FindDeviceCode)
		device.POST("/name", FindDeviceName)
		device.POST("/value", DeviceValue)
	}

	template := e.Group("/template") //需要验证token
	template.Use(jwt.JWTAuth())
	{
		template.GET("/all", AllTem)
		template.POST("/add", AddTem)
		template.POST("/update", UpdateTem)
		template.POST("/remove", RemoveTem)
		template.POST("/key", FindTem)
	}
	//数据库管理模块的请求
	dbMgr := e.Group("/dbMgr")  //需要验证token
	dbMgr.Use(jwt.JWTAuth())
	{
		dbMgr.GET("/all",AlldbMgr)	//查询出所有的数据库数据
		dbMgr.POST("/add",AddDbMgr)	//添加一条记录
		dbMgr.POST("/update",UpdateDbMgr)	//更新一条记录
		dbMgr.POST("/delete",DeleteDbMgr)	//删除一条记录
		dbMgr.POST("/find",FindServerIp)	//查找当前ip的数据库
		dbMgr.POST("/ping", TestPing)		//数据库测试连接
	}

	e.GET("/alarm/all", AllAlarm)

	e.GET("/alarmStg/all", AllStg)
	e.POST("/alarmStg/add", NewStg)
	e.POST("/alarmStg/update", UpdateStg)
	e.POST("/alarmStg/remove", RemoveStg)
	e.POST("/alarmStg/key", FindStg)
	return e
}

//header添加token验证后，默认的跨域头失效，重写一个
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, User,Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
