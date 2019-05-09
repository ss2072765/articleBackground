package controllers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
	"project01/models"
)

func init1(){
	con,e:=redis.Dial("tcp","127.0.0.1:6379")
	if e!=nil{
		beego.Error(e)
		return
	}
	defer con.Close()
	//con.Send("set","key1","hello world222")
	//con.Send("get","key1")
	//con.Flush()
	//ret,e:=con.Receive()
	//ret,e:=con.Do("get","key1")
	ret,e:=con.Do("hmget","user","Id","Name","Pwd")
	if e!=nil{
		beego.Error(e)
		return
	}
	slice,e:=redis.Values(ret,e)
	if e!=nil{
		beego.Error(e)
		return
	}
	var user models.User
	redis.ScanSlice(slice,&user)
	beego.Info(user)

}
