package main

import (
	_ "project01/routers"
	"github.com/astaxie/beego"
	_ "project01/models"
	"github.com/astaxie/beego/orm"
)

func main() {
	orm.Debug=true
	beego.AddFuncMap("prePageFunc",prePage)
	beego.AddFuncMap("nextPageFunc",nextPage)
	beego.AddFuncMap("addIndex",addIndex)
	beego.Run()
}

func prePage(pageNumber int) int{
	if pageNumber<=1{
		return 1
	}
	return pageNumber-1
}

func nextPage(pageNumber int,pageCount float64) int{
	if pageNumber>=int(pageCount){
		return pageNumber
	}
	return pageNumber+1
}

func addIndex(index int) int{
	return index+1
}