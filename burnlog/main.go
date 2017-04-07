// main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
)

const (
	serviceLabel = ":func"
	addonLabel   = ":addon"
)

func main() {
	routeRule := fmt.Sprintf("/%s/%s", serviceLabel, addonLabel)
	routeRule2 := fmt.Sprintf("/%s", serviceLabel)
	fmt.Println(routeRule)
	beego.Router(routeRule, &mainController{})
	beego.Router(routeRule2, &mainController{})
	beego.Run()
}

type mainController struct {
	beego.Controller
}

func (c *mainController) Get() {
	operation := c.Ctx.Input.Param(serviceLabel)
	fmt.Println(operation)
	//	addon := c.Ctx.Input.Param(addonLabel)
	w := c.Ctx.ResponseWriter
	//	r := c.Ctx.Request
	switch operation {
	case "article":
		showTempletResponse(w, operation)
	case "user":
		showTempletResponse(w, operation)
	case "material":
		showTempletResponse(w, operation)
	case "comment":
		showTempletResponse(w, operation)
	default:
		httpNotFound(w)
	}
}

func showTempletResponse(w http.ResponseWriter, templete string) {
	text := fmt.Sprintf("The API is :	%s", templete)
	fmt.Fprintln(w, text)
}

func httpNotFound(w http.ResponseWriter) {
	fmt.Fprintln(w, "Sorry, Our blog is not open yet!")
}
