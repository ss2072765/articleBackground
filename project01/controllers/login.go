package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"project01/models"
	"encoding/base64"
	"github.com/gomodule/redigo/redis"
)

type LoginController struct{
	beego.Controller
}

func(this *LoginController)Login(){
	name:=this.Ctx.GetCookie("name")
	base64,_:=base64.StdEncoding.DecodeString(name)
	remember:=this.Ctx.GetCookie("remember")
	this.Data["name"]=string(base64)
	this.Data["remember"]=remember
	this.TplName="login.html"
}

func(this *LoginController)HandleLogin(){
	name:=this.GetString("userName")
	pwd:=this.GetString("password")
	remember:=this.GetString("remember")
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
		this.TplName="login.html"
		return
	}
	ormer:=orm.NewOrm()
	var user models.User
	e:=ormer.Raw("select * from user where name=? and pwd=?",name,pwd).QueryRow(&user)
	if e!=nil{
		beego.Error(e)
		this.Data["nameErr"]="用户名不存在或密码错误"
		this.TplName="login.html"
		return
	}
	if remember=="on"{
		base64:=base64.StdEncoding.EncodeToString([]byte(name))
		this.Ctx.SetCookie("name",base64,60)
		this.Ctx.SetCookie("remember","checked",60)
	}else{
		this.Ctx.SetCookie("name",name,-1)
		this.Ctx.SetCookie("remember",remember,-1)
	}
	this.SetSession("user",user)
	//重定向至文章列表页
	this.Redirect("/toHomePage",302)
}

func(this *LoginController) Logout(){
	con,e:=redis.Dial("tcp","127.0.0.1:6379")
	if e!=nil{
		beego.Error(e)
		return
	}
	_,e=con.Do("del","articleTypes")
	if e!=nil{
		beego.Error(e)
		return
	}
	this.DelSession("user")
	this.Redirect("/login",302)
}