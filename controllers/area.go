package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
)

type AreaController struct {
	beego.Controller
}

/*
type AreaResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}
map[string]interface{}
*/

//将封装好的返回结构 变成json返回给前段
func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//  /api/v1.0/areas [get]
func (this *AreaController) GetAreaInfo() {
	beego.Info("==========/api/v1.0/area get succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//1 从缓存中redis读数据

	//2 如果redis有 之前的json字符串数据那么直接返回给前段

	//3 如果redis没有之前的json字符串数据， 从mysql查
	o := orm.NewOrm()

	//得到查到的areas数据
	var areas []models.Area //[{aid,aname},{aid,aname},{aid,aname}]

	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		//返回错误信息给前端
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)

		return
	}
	if num == 0 {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)

		return
	}

	//succ
	resp["data"] = areas

	//将封装好的返回结构体map 发送给前段
	return
}
