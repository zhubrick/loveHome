package routers

import (
	"github.com/astaxie/beego"
	"loveHome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	//请求地域信息
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreaInfo")

	//session
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:GetSessionName;delete:DelSessionName")
	//登陆
	/*
		method: POST
		api/v1.0/sessions
	*/
	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	// api/v1.0/users [post] //注册
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")

	//house/index  房屋首页列表
	beego.Router("/api/v1.0/houses/index", &controllers.HousesIndexController{}, "get:GetHousesIndex")

}
