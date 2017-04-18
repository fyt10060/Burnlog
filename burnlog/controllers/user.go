// user
package controllers

import (
	"Burnlog/burnlog/model"
	"Burnlog/burnlog/service"

	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type UserController struct {
	beego.Controller
}

const (
	signup = "signup"
	signin = "signin"
	users  = "/"
	del    = "delete"
	// router
	UserParam = ":user"
	// sign in & sign up
	paramUserName = "username"
	paramPassword = "password"
	paramEmail    = "email"
	// normal const
	isAdmainCount = "authority"
)

func (this *UserController) Get() {
	addon := this.Ctx.Input.Param(UserParam)
	switch addon {
	case signin:
		signIn(this.Ctx)
	case signup:
		signUp(this.Ctx)
	case users:
		fmt.Fprintln(this.Ctx.ResponseWriter, "Your are seeing users page")
	default:
		break
	}
	//	fmt.Fprintf(this.Ctx.ResponseWriter, "action is :%s\n", addon)
	fmt.Fprint(this.Ctx.ResponseWriter)
}

func (this *UserController) Post() {
	addon := this.Ctx.Input.Param(UserParam)
	switch addon {
	case signin:
		signIn(this.Ctx)
	case signup:
		signUp(this.Ctx)
	default:
		break
	}
}

func signIn(ctx *context.Context) {
	r := ctx.Request
	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	name, password, email, result := getUserInfo(r)
	if result != nil {
		w.Write(result)
		return
	}
	if name == "" || email == "" {
		result = model.ParseResult(model.ErrNameEmpty, nil)
		w.Write(result)
		return
	}

	msg, userInfo := service.Login(email, password)
	result = model.ParseResult(msg, userInfo)
	w.Write(result)
}

func signUp(ctx *context.Context) {
	r := ctx.Request
	r.ParseForm()

	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	name, password, email, result := getUserInfo(r)

	if result != nil {
		w.Write(result)
		return
	}
	if name == "" && email == "" {
		result := model.ParseResult(model.ErrNameEmpty, nil)
		w.Write(result)
		return
	}
	newUser := model.CreateNewUser(name, email)
	msg := service.AddNewUser(newUser, password)
	result = model.ParseResult(msg, nil)
	w.Write(result)
}

func getUserInfo(r *http.Request) (name, password, email string, passwordError []byte) {

	pasw := r.Form.Get(paramPassword)
	pasLen := len(pasw)
	//	fmt.Printf("pass len is %d", pasLen)
	if pasLen < 8 || pasLen > 20 {
		result := model.ParseResult(model.ErrPassIllegal, nil)
		return "", "", "", result
	}
	accName := r.Form.Get(paramUserName)
	email = r.Form.Get(paramEmail)

	return accName, pasw, email, nil
}
