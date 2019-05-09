package routers

import (
	"project01/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/*",beego.BeforeExec,filter)
    //beego.Router("/", &controllers.MainController{})
    //跳转至注册页面
    beego.Router("/index",&controllers.RegisterController{},"get:Index")
    //注册接口
    beego.Router("/register",&controllers.RegisterController{},"post:Register")
    //跳转至登陆页面
	beego.Router("/login",&controllers.LoginController{},"get:Login")
	//登陆接口
	beego.Router("/handleLogin",&controllers.LoginController{},"post:HandleLogin")
	//跳转至首页即文章列表页
	beego.Router("/toHomePage",&controllers.ArticleController{},"get:ToHomePage")
	//跳转至文章详情页面
	beego.Router("/toContent",&controllers.ArticleController{},"get:ToContent")
	//跳转至文章添加页面
	beego.Router("/toAdd",&controllers.ArticleController{},"get:ToAdd")
	//跳转至文章修改页面
	beego.Router("/toUpdate",&controllers.ArticleController{},"get:ToUpdate")
	//文章修改接口
	beego.Router("/update",&controllers.ArticleController{},"post:Update")
	//文章删除接口
	beego.Router("/delete",&controllers.ArticleController{},"get:Delete")
	//跳转至文章“类型”添加页面
	beego.Router("/toAddType",&controllers.ArticleController{},"get:ToAddType")
	//添加文章类型接口
	beego.Router("/addType",&controllers.ArticleController{},"post:AddType")
	//删除文章类型接口
	beego.Router("/deleteType",&controllers.ArticleController{},"get:DeleteType")
	//添加文章接口
	beego.Router("/add",&controllers.ArticleController{},"post:Add")
	//退出接口
	beego.Router("/logout",&controllers.LoginController{},"get:Logout")


}

func filter(context *context.Context)  {
	user:=context.Input.Session("user")
	url:=context.Request.URL.Path
	if url!="/index"&&url!="/register"&&url!="/login"&&url!="/handleLogin"{
		if user==nil{
			context.Redirect(302,"/login")
			return
		}
	}

}