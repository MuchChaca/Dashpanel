package main

import (
	"github.com/MuchChaca/Dashpanel/src/controller"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
)

func main() {
	app := iris.New()

	// serve our app in public, public folder
	// contains the client-side vue.js application,
	// no need for any server-side template here,
	// actually if you're going to use vue without any
	// back-end services, you can just stop after this line and start the server.
	app.StaticWeb("/", "../public/dist")

	// configure the http sessions.
	// sess := sessions.New(sessions.Config{
	// 	Cookie: "iris_session",
	// })

	// configure the websocket server.
	ws := websocket.New(websocket.Config{})

	// // DASH
	// create a sub router an register the client-side library for the iris websockets,
	// you could skip it but iris websockets supports socket.io-like API.
	dashRouter := app.Party("/dash")
	// http://localhost:8080/todos/iris-ws.js
	// serve the javascript client library to communicate with
	// the iris high level websocket event system.
	dashRouter.Any("/iris-ws.js", websocket.ClientHandler())

	//create our mvc app targeted to /dash relative sub path.
	dashApp := mvc.New(dashRouter)

	// any dependencies bindings here . . .
	dashApp.Register(
		// dash.NewMemoryService(),
		// sess.Start,
		ws.Upgrade,
	)
	// controllers registration here
	dashApp.Handle(new(controller.DashController))

	// start the web server at http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithoutVersionChecker)
}
