// main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
)

const (
	serviceLabel = "/:func"
	addonLabel   = "/:addon"
)

func main() {
	routeRule := fmt.Sprintf("/%s/%s", serviceLabel, addonLabel)
	beego.Router(routeRule, &mainController{})
	beego.Run()
}

type mainController struct {
	beego.Controller
}

func (c *mainController) Get() {
	operation := c.Ctx.Input.Param(serviceLabel)
	addon := c.Ctx.Input.Param(addonLabel)
	switch operation {
	case "article":
	default:
	}
}

func httpNotFound(w http.ResponseWriter) {
	fmt.Fprintln("")
}
