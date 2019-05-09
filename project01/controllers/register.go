package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RegisterController struct{
	beego.Controller
}

func(this *RegisterController) Index(){
	 this.TplName="register.html"
}

func(this *RegisterController) Register(){
	var name=this.GetString("userName")
	var pwd=this.GetString("password")
	var b=false
	if name==""{
		this.Data["nameErr"]="用户名不能为空！"
		b=true
	}
	if pwd==""{
		this.Data["pwdErr"]="密码不能为空！"
		b=true
	}
	if b{
		this.TplName="register.html"
		return
	}
	ormer:=orm.NewOrm()
	result,e:=ormer.Raw("insert into user values(null,?,?)",name,pwd).Exec()
	if e!=nil{
		beego.Error(e)
		return
	}
	beego.Info(result.LastInsertId())
	this.TplName="login.html"
}

