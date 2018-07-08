package controllers

import (
	"github.com/kataras/iris/_examples/mvc/login/services"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

type HomeController struct {
	Ctx iris.Context
	Service services.UserService
	Session *sessions.Session
}

func (c *HomeController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

func (c *HomeController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *HomeController) logout() {
	c.Session.Destroy()
}

func (c *HomeController) Get() mvc.Result {
	if !c.isLoggedIn() {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/user/login"}
	}

	u, found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		c.logout()
		return c.Get()
	}

	return mvc.View{
		Name: "pages/index.html",
		Data: iris.Map{
			"Title": "Welcome " + u.Username,
			"User":  u,
		},
	}
}

func (c *HomeController) GetMe() mvc.Result {

	// ถ้าไม่ได้ login ให้ไป login
	if !c.isLoggedIn() {
		return mvc.Response{Path: "/user/login"}
	}

	// ถ้าไม่มีข้อมูล User ID ของฉันให้ logout
	u, found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		c.logout()
		return c.GetMe()
	}

	return mvc.View{
		Name: "pages/me.html",
		Data: iris.Map{
			"Title": "Welcome " + u.Username,
			"User":  u,
		},
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *HomeController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}

	c.Ctx.Redirect("/user/login")
}
