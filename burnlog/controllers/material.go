// material
package controllers

import (
	//	"fmt"
	"net/http"

	//	"Burnlog/burnlog/service"

	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/context"
)

const (
	// action name
	MaterialParam = ":material"
	actMeList     = "list"
	actMeUpload   = "upload"
	actMeDelete   = "delete"
	actMeGetId    = "get"
	// param name
)

type MaterialController struct {
	beego.Controller
}

func (this *MaterialController) Get() {
	c := this.Ctx
	r := c.Request
	w := c.ResponseWriter

	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	action := c.Input.Param(MaterialParam)
	switch action {
	case actMeDelete:
		materialDelete(r, w)
	case actMeList:
		materialList(r, w)
	case actMeGetId:
		getMaterialWithId(r, w)
	case actMeUpload:
		materialUpload(r, w)
	default:
		break
	}
}

func getMaterialWithId(r *http.Request, w http.ResponseWriter) {

}

func materialDelete(r *http.Request, w http.ResponseWriter) {

}

func materialUpload(r *http.Request, w http.ResponseWriter) {

}

func materialList(r *http.Request, w http.ResponseWriter) {

}
