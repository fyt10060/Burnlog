// article
package controllers

import (
	"fmt"
	"net/http"

	"Burnlog/burnlog/model"
	//	"Burnlog/burnlog/service"

	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/context"
)

const (
	ArticleParam = ":action"

	// action name
	articleList   = "list"
	createArticle = "create"
	updateArticle = "update"
	deleteArticle = "delete"

	// param name
	paramArticleId = "article_id"

	// test
	testArticleId = "testarticle001"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) Get() {
	c := this.Ctx
	r := c.Request
	r.ParseForm()
	w := c.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	action := c.Input.Param(ArticleParam)
	switch action {
	case articleList:
		List(r, w)
	case createArticle:
		//		Create(r, w)
	case updateArticle:
		//		Update(r, w)
	case deleteArticle:
		Delete(r, w)
	default:
		fmt.Fprintln(c.ResponseWriter, "Your are seeing the article page!")
	}
}

// Get article list
func List(r *http.Request, w http.ResponseWriter) {
	var list model.ArticleList
	list.List = append(list.List, model.Article{
		Title:      "我们在大草原的湖边",
		Detail:     "等候鸟飞回来",
		CreateTime: 1491969188,
		Author:     "flamel",
	})
	list.List = append(list.List, model.Article{
		Title:      "等我们都长大了",
		Detail:     "就生一个娃娃，他会自己长大远去，而我们也各自老去，我给你写信",
		CreateTime: 1491969100,
		Author:     "gogogo",
	})
	//	result := model.ParseError(0, "success", list)
	//	w.Write(result)
}

// Create article
type CreateResult struct {
	CreateTime int64  `json:"create_time"`
	ArticleID  string `json:"article_id"`
}

//func Create(r *http.Request, w http.ResponseWriter) {

//	dataCorrect, result := checkArticle(r)
//	if !dataCorrect {
//		w.Write(result)
//		return
//	}

//	nowTime := service.GetTime()
//	data := CreateResult{
//		CreateTime: nowTime,
//		ArticleID:  testArticleId,
//	}
//	//	result = model.ParseError(0, "success", data)
//	w.Write(result)
//}

// Make sure article upload not wrong
//func checkArticle(r *http.Request) (bool, []byte) {
//	title := r.Form.Get("title")
//	if len(title) == 0 {
////		result := model.ParseError(20001, "title is empty", nil)
//		return false, result
//	}
//	detail := r.Form.Get("detail")
//	if len(detail) == 0 {
////		result := model.ParseError(20002, "text is empty", nil)
//		return false, result
//	}
//	return true, nil
//}

// Delete article
func Delete(r *http.Request, w http.ResponseWriter) {
	result := checkArticleId(r)
	if result != nil {
		w.Write(result)
		return
	}
	//	result = model.ParseError(0, "success", nil)
	w.Write(result)
}

// Make sure id exist
func checkArticleId(r *http.Request) []byte {
	//	id := r.Form.Get(paramArticleId)
	//	if id != testArticleId {
	//		result := model.ParseError(20003, "article id not found", nil)
	//		return result
	//	}
	return nil
}

// Update article
//func Update(r *http.Request, w http.ResponseWriter) {
//	dataCorrect, result := checkArticle(r)
//	if !dataCorrect {
//		w.Write(result)
//		return
//	}
//	result = checkArticleId(r)
//	if result != nil {
//		w.Write(result)
//		return
//	}
//	//	result = model.ParseError(0, "success", nil)
//	w.Write(result)
//}
