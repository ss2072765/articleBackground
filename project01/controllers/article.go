package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"project01/models"
	"math"
	"github.com/gomodule/redigo/redis"

	"encoding/json"
	"encoding/base64"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ToHomePage() {
	//this.ServeJSON()
	//user:=this.GetSession("user")
	//if userName==nil{
	//	this.Redirect("/login",302)
	//	return
	//}
	//this.Data["userName"]=user.(models.User).Name
	typeId := this.GetString("select")
	pageNumber, e := this.GetInt("pageNumber")
	if e != nil {
		beego.Error(e)
		pageNumber = 1
	}
	ormer := orm.NewOrm()
	var count int
	e = ormer.Raw("select count(id) from article where if(?='',1=1,article_type_id=?)",typeId,typeId).QueryRow(&count)
	if e != nil {
		beego.Error(e)
		return
	}
	//beego.Info(count)
	var perCount = 3                                              //设置每页记录数
	var pageCount = math.Ceil(float64(count) / float64(perCount)) //向上取整求总页数
	var articles []orm.Params
	_, e = ormer.Raw("select a.id id,a.title title,a.content content,a.read_count readCount,a.create_at createAt,b.type_name typeName "+
		"from article a left join article_type b on a.article_type_id=b.id "+
		"where if(?=0,1=1,b.id=?)"+
		"order by a.create_at desc limit?,?",
		typeId, typeId, (pageNumber-1)*int(perCount), perCount).Values(&articles)
	if e != nil {
		beego.Error(e)
		return
	}

	var articleTypes []orm.Params
	con,e:=redis.Dial("tcp","127.0.0.1:6379")
	if e!=nil{
		beego.Error(e)
		return
	}
	defer con.Close()
	inter,e:=con.Do("get","articleTypes")
	if e!=nil{
		beego.Error(e)
		return
	}
	if inter==nil{
		_, e = ormer.Raw("select id Id,type_name TypeName from article_type order by Id desc").Values(&articleTypes)
		//var buffer bytes.Buffer
		//encoder:=gob.NewEncoder(&buffer)
		//e=encoder.Encode(articleTypes)
		//if e!=nil{
		//	beego.Error(e)
		//	return
		//}
		bytes,e:=json.Marshal(articleTypes)
		if e!=nil{
			beego.Error(e)
			return
		}
		_,e=con.Do("set","articleTypes",bytes)
		if e!=nil{
			beego.Error(e)
			return
		}
		beego.Info("mysql")
	}else{
		result,e:=redis.Bytes(inter,e)
		if e!=nil{
			beego.Error(e)
			return
		}
		//decoder:=gob.NewDecoder(bytes.NewBuffer(result))
		//decoder.Decode(&articleTypes)
		json.Unmarshal(result,&articleTypes)
		beego.Info("redis")
	}

	//var articleTypes []orm.Params
	//con,e:=redis.Dial("tcp","127.0.0.1:6379")
	//if e!=nil{
	//	beego.Error(e)
	//	return
	//}
	//defer con.Close()
	//inter,e:=con.Do("get","articleTypes")
	//if e!=nil{
	//	beego.Error(e)
	//	return
	//}
	//if inter==nil{
	//	_, e = ormer.Raw("select id Id,type_name TypeName from article_type order by Id desc").Values(&articleTypes)
	//	var buffer bytes.Buffer
	//	encoder:=gob.NewEncoder(&buffer)
	//	encoder.Encode(articleTypes)
	//	_,e:=con.Do("set","articleTypes",buffer.Bytes())
	//	if e!=nil{
	//		beego.Error(e)
	//		return
	//	}
	//	beego.Info("mysql")
	//}else{
	//	result,e:=redis.Bytes(inter,e)
	//	if e!=nil{
	//		beego.Error(e)
	//		return
	//	}
	//	encoder:=gob.NewDecoder(bytes.NewBuffer(result))
	//	encoder.Decode(&articleTypes)
	//	beego.Info("redis")
	//}




	if e != nil {
		beego.Error(e)
		return
	}
	if typeId != "" {
		for i := 0; i < len(articleTypes); i++ {
			if articleTypes[i]["Id"].(string) == typeId {
				articleTypes[i]["Select"] = "selected"
				this.Data["typeId"]=typeId
			} else {
				articleTypes[i]["Select"] = ""
			}
		}
	}
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageNumber"] = pageNumber
	this.Data["articles"] = articles
	this.Data["articleTypes"] = articleTypes
	this.TplName = "index.html"
}

func (this *ArticleController) ToContent() {
	id, e := this.GetInt("id")
	if e != nil {
		beego.Error(e)
		this.Redirect("/toHomePage", 302)
		return
	}
	ormer := orm.NewOrm()
	var article models.Article
	e = ormer.Raw("select * from article where id=?", id).QueryRow(&article)
	if e != nil {
		beego.Error(e)
		return
	}
	article.ReadCount += 1
	result, e := ormer.Raw("update article set read_count=? where id=?", article.ReadCount, id).Exec()
	if e != nil {
		beego.Error(e)
		return
	}
	beego.Info(result.LastInsertId())
	userId:=this.GetSession("user").(models.User).Id
	_,e=ormer.Raw("insert into article_users values(null,?,?)",id,userId).Exec()
	if e!=nil{
		beego.Error(e)
		return
	}
	var userNames []orm.Params
	_,e=ormer.Raw("select distinct a.* from (select concat(b.name,' | ') userName from article_users a left join user b on a.user_id = b.id where a.article_id=? order by a.id desc) a",id).Values(&userNames)
	this.Data["article"] = article
	this.Data["userNames"] = userNames
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

func (this *ArticleController) ToAdd() {
	ormer := orm.NewOrm()
	var articleTypes []models.ArticleType
	_, e := ormer.Raw("select * from article_type order by id desc").QueryRows(&articleTypes)
	if e != nil {
		beego.Error(e)
		return
	}
	this.Data["articleTypes"] = articleTypes
	this.TplName = "add.html"
}

func (this *ArticleController) ToUpdate() {
	id, e := this.GetInt("id")
	if e != nil {
		beego.Error(e)
		this.Data["errmsg"] = "文章Id解析错误"
		this.TplName = "update.html"
	}
	ormer := orm.NewOrm()
	var article models.Article
	e = ormer.Raw("select * from article where id=?", id).QueryRow(&article)
	if e != nil {
		beego.Error(e)
		return
	}
	this.Data["article"] = article
	this.TplName = "update.html"
}

func (this *ArticleController) Update() {
	id := this.GetString("id")
	title := this.GetString("articleName")
	content := this.GetString("content")
	file, head, e := this.GetFile("uploadname")
	defer file.Close()
	if e != nil {
		beego.Error(e)
		//this.Data["errmsg"]="文件上传失败"
		this.Redirect("/toUpdate?id="+id, 302)
		return
	}
	if head.Size > 5000000 {
		beego.Error("文件上传过大")
		this.Redirect("/toUpdate?id="+id, 302)
		return
	}
	var filePath string
	if head.Filename == "" && head.Size == 0 {
		filePath = ""
	} else {
		ext := path.Ext(head.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			beego.Error("文件格式不正确")
			this.Redirect("/toUpdate?id="+id, 302)
			return
		}
		nowFormat := time.Now().Format("2006010215040506")
		filePath = "/static/img" + nowFormat + ext
		this.SaveToFile("uploadname", "."+filePath)
	}

	ormer := orm.NewOrm()
	result, e := ormer.Raw("update article set title=?,content=?,img=? where id=?", title, content, filePath, id).Exec()
	if e != nil {
		beego.Error(e)
		return
	}
	beego.Info(result.LastInsertId())
	this.Redirect("/toHomePage", 302)
}

func (this *ArticleController) Delete() {
	id := this.GetString("id")
	ormer := orm.NewOrm()
	result, e := ormer.Raw("delete from article where id=?", id).Exec()
	if e != nil {
		beego.Error(e)
		return
	}
	beego.Info(result.LastInsertId())
	this.Redirect("/toHomePage", 302)
}

func (this *ArticleController) ToAddType() {
	ormer := orm.NewOrm()
	var articleTypes []models.ArticleType
	_, e := ormer.Raw("select * from article_type order by id desc").QueryRows(&articleTypes)
	if e != nil {
		beego.Error(e)
		return
	}
	this.Data["articleTypes"] = articleTypes
	this.TplName = "addType.html"
}

func (this *ArticleController) AddType() {
	typeName := this.GetString("typeName")
	ormer := orm.NewOrm()
	_, e := ormer.Raw("insert into article_type values(null,?)", typeName).Exec()
	if e != nil {
		beego.Error(e)
		return
	}
	this.Redirect("/toAddType", 302)
}

func (this *ArticleController) DeleteType() {
	id := this.GetString("id")
	ormer := orm.NewOrm()
	_, e := ormer.Raw("delete from article_type where id=?", id).Exec()
	if e != nil {
		beego.Error(e)
		return
	}
	this.Redirect("/toAddType", 302)
}

func (this *ArticleController) Add() {
	articleTypeId, e := this.GetInt("select")
	if e != nil {
		beego.Error(e)
		this.Data["errmsg"] = "文章类型ID有误"
		this.TplName = "add.html"
		return
	}
	var title = this.GetString("articleName")
	if title == "" {
		this.Data["errmsg"] = "文章标题不能为空"
		this.TplName = "add.html"
		return
	}
	var content = this.GetString("content")
	if content == "" {
		this.Data["errmsg"] = "文章内容不能为空"
		this.TplName = "add.html"
		return
	}
	file, head, e := this.GetFile("uploadname")
	defer file.Close()
	if e != nil {
		beego.Error(e)
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}
	if head.Size > 5000000 {
		this.Data["errmsg"] = "图片上传过大"
		this.TplName = "add.html"
		return
	}
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "上传文件格式错误"
		this.TplName = "add.html"
		return
	}
	newFileName := time.Now().Format("2006010215040506")
	path := "/static/img/" + newFileName + ext
	this.SaveToFile("uploadname", "."+path)
	//now:=time.Now()
	//beego.Info(now)
	ormer := orm.NewOrm()
	result, e := ormer.Raw(
		"insert into article(title,content,img,create_at,article_type_id) values(?,?,?,now(),?)",
		title, content, path, articleTypeId).Exec()
	if e != nil {
		beego.Error(e)
		this.Data["errmsg"] = "插入数据库失败"
		this.TplName = "add.html"
		return
	}
	beego.Info(result)
	this.Redirect("/toHomePage", 302)

}
