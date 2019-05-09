package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)


type User struct{
	Id int
	Name string
	Pwd string
	Article	[]*Article `orm:"reverse(many)"`
}

type Article struct{
	Id int `orm:"pk;auto"` //主键
	Title string `orm:"unique;size(50)"` //文章标题
	Content string `orm:"size(500)"`  //文章内容
	ReadCount int `orm:"default(0)"`  //阅读量
	Img string	`orm:"null"`  //图片路径
	CreateAt time.Time `orm:"type(datatime);auto_now_add"`  //创建时间
	ArticleType *ArticleType `orm:"rel(fk)"`
	User []*User `orm:"rel(m2m)"`
}

type ArticleType struct{
	Id int
	TypeName string
	Article []*Article `orm:"reverse(many)"`
}
func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/project01")
	//var (
	//	a Article
	//	b ArticleType
	//	c User
	//	)
	//orm.RegisterModel(&a,&b,&c)
	//orm.RunSyncdb("default",false,true)
}