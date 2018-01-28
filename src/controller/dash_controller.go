package controller

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/MuchChaca/Dashpanel/src/model/dash"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

// DashController is the controller for "dash"
type DashController struct {
	Service dash.Service
	Process dash.Process

	Session *sessions.Session
}

// BeforeActiviation called once before the server ran, and before
// the routes and dependencies binded.
// You can bind custom things to the cont`roller, add new methods, add middleware,
// add dependencies to the struct or the method(s) and more.
func (c *DashController) BeforeActiviation(bfA mvc.BeforeActivation) {
	// this could be binded to a controller's function input argument
	// if any, or struct field if any:
	bfA.Dependencies().Add(func(ctx iris.Context) (items []dash.Process) {
		ctx.ReadJSON(&items)
		return
	})
}

// Get handles the GET: /dash route.
func (c *DashController) Get() iris.Map {
	return c.LoadProcess()
}

// LoadProcess loads the processes
func (c *DashController) LoadProcess() iris.Map {
	topCmd := exec.Command("pstree", "-ptc")
	var topOut bytes.Buffer
	topCmd.Stdout = &topOut

	errTop := topCmd.Run()
	if errTop != nil {
		log.Fatalln("WTF DUDE !!?", errTop)
	}

	outputResult := fmt.Sprintf("%q", topOut.String())

	rgex, _ := regexp.Compile("-([a-zA-Z]|\\d)+\\(\\d+\\)")
	rgexResult := rgex.FindAllString(outputResult, -1)

	var allProcess map[string]interface{}
	allProcess = make(map[string]interface{})

	httpd := dash.TreeProcess{
		Tree: make([]dash.Processus, 0),
	}

	docker := dash.TreeProcess{
		Tree: make([]dash.Processus, 0),
	}

	// otherProc := dash.TreeProcess{
	// 	Tree: make([]dash.Processus, 2),
	// }

	// i := 0
	for _, str := range rgexResult {
		str = strings.Replace(str, "-", "", -1)
		rgex, _ = regexp.Compile("^([a-zA-Z])([a-zA-Z]|\\d)+")
		processusName := rgex.FindString(str)

		rgex, _ = regexp.Compile("(\\d+)")
		processusPID, errConv := strconv.Atoi(rgex.FindString(str))
		if errConv != nil {
			log.Panicln("Error in convertion PID", errConv)
		}
		// index := strconv.Itoa(i)
		aProcess := dash.Processus{
			Name:    processusName,
			PID:     processusPID,
			Version: "N/A",
			Status:  true,
		}

		// allProcess[processusName] = append(allProcess[processusName], aProcess)

		/////////
		switch aProcess.Name {
		case "dockerd":
			docker.Tree = append(docker.Tree, aProcess)
		case "httpd":
			httpd.Tree = append(httpd.Tree, aProcess)
		default:
			// otherProc.Tree = append(otherProc.Tree, aProcess)
		}
		/////////

		// i++
	}
	// fmt.Println(rgexResult)

	allProcess["Docker"] = docker
	allProcess["Apache"] = httpd
	// allProcess["Others"] = otherProc

	return allProcess
}
