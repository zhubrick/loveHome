package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
)

type UserController struct {
	beego.Controller
}

//将封装好的返回结构 变成json返回给前段
func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//  /api/v1.0/users [post]
/*

{
	mobile: "123",
	password: "123",
	sms_code: "123"
}
*/
func (this *UserController) Reg() {
	beego.Info("==========/api/v1.0/users post succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	var regRequestMap = make(map[string]interface{})

	//1 得到客户端请求的json数据 post数据
	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestMap)

	beego.Info("mobile = ", regRequestMap["mobile"])
	beego.Info("password = ", regRequestMap["password"])
	beego.Info("sms_code = ", regRequestMap["sms_code"])

	//2 判断数据的合法性
	if regRequestMap["mobile"] == "" || regRequestMap["password"] == "" || regRequestMap["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//3 将数据存入到mysql数据库 user
	user := models.User{}
	user.Mobile = regRequestMap["mobile"].(string)
	//应该将password进行md5，SHA246,SHA1
	user.Password_hash = regRequestMap["password"].(string)
	user.Name = regRequestMap["mobile"].(string)

	o := orm.NewOrm()

	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("insert error = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	beego.Info("reg succ !!! user id = ", id)

	//4 将当前的用户的信息存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)

	return
}

//登陆
/*
	method: POST
	api/v1.0/sessions

	{
		mobile: "133",
		password: "itcast"
	}
*/
func (this *UserController) Login() {
	beego.Info("==========/api/v1.0/sessions login succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	var loginRequestMap = make(map[string]interface{})

	//1 得到客户端请求的json数据 post数据
	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap)

	beego.Info("mobile = ", loginRequestMap["mobile"])
	beego.Info("password = ", loginRequestMap["password"])

	//2 判断数据的合法性
	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//3 查询数据库得到user
	var user models.User

	o := orm.NewOrm()
	//select password from user where user.name = name
	qs := o.QueryTable("user")
	if err := qs.Filter("mobile", loginRequestMap["mobile"]).One(&user); err != nil {
		//查询失败
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}

	//4 对比密码
	if user.Password_hash != loginRequestMap["password"].(string) {
		resp["errno"] = models.RECODE_PWDERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PWDERR)
		return
	}

	beego.Info("==== login succ!!! === user.name = ", user.Name)

	//5 将当前的用户的信息存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}
