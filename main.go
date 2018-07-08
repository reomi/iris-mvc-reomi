package main

import (
	"time"
	"github.com/kataras/iris/_examples/mvc/login/repositories"

	"iris-reomi/datasource"
	"iris-reomi/services"
	"iris-reomi/web/middleware"
	"iris-reomi/web/controllers"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	app.StaticWeb("/public", "./web/public")
	app.StaticWeb("assets", "./web/views/themes/assan-admin/assets")

	// Load the template files.
	tmpl := iris.HTML("./web/views/themes/assan-admin", ".html").
		Layout("layouts/layout-main.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("layouts/layout.html")
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("pages/error.html")
	})

	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}
	repo := repositories.NewUserRepository(db)
	userService := services.NewUserService(repo)

	// API "/users" based mvc application.
	users := mvc.New(app.Party("/users"))
	// กำหนด BasicAuth ให้กับ api users
	users.Router.Use(middleware.BasicAuth)

	// Bind the "userService" to the UserController's Service (interface) field.
	users.Register(userService)
	users.Handle(new(controllers.UsersController))

	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})

	// สร้าง Controllers "User"
	user := mvc.New(app.Party("/user").Layout("layouts/layout.html"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))

	// สร้าง Controllers "Home"
	home := mvc.New(app.Party("/"))
	home.Register(
		userService,
		sessManager.Start,
	)
	home.Handle(new(controllers.HomeController))

	// http://localhost:8080/noexist
	// and all controller's methods like
	// http://localhost:8080/users/1
	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
