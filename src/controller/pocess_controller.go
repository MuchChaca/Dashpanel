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
	"github.com/fatih/color"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

// ProcessController is the controller for "dash"
type ProcessController struct {
	Service dash.Service
	Process dash.Process

	Session *sessions.Session
}

// BeforeActiviation called once before the server ran, and before
// the routes and dependencies binded.
// You can bind custom things to the cont`roller, add new methods, add middleware,
// add dependencies to the struct or the method(s) and more.
func (c *ProcessController) BeforeActiviation(bfA mvc.BeforeActivation) {
	// this could be binded to a controller's function input argument
	// if any, or struct field if any:
	bfA.Dependencies().Add(func(ctx iris.Context) (items []dash.Process) {
		ctx.ReadJSON(&items)
		return
	})
}

// Get handles the GET: /dash route.
func (c *ProcessController) Get() iris.Map {
	return c.LoadProcess()
}

// LoadProcess loads the processes
func (c *ProcessController) LoadProcess() iris.Map {
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

	// Sadly, one for each process that we want a group of - names are revelant...
	httpd := dash.TreeProcess{
		Tree: make([]dash.Processus, 0),
	}

	docker := dash.TreeProcess{
		Tree: make([]dash.Processus, 0),
	}

	mysqld := dash.TreeProcess{
		Tree: make([]dash.Processus, 0),
	}

	// // To load all the other processes in as "Others"
	// otherProc := dash.TreeProcess{
	// 	Tree: make([]dash.Processus, 2),
	// }

	for _, str := range rgexResult {
		str = strings.Replace(str, "-", "", -1)
		rgex, _ = regexp.Compile("^([a-zA-Z])([a-zA-Z]|\\d)+")
		processusName := rgex.FindString(str)

		rgex, _ = regexp.Compile("(\\d+)")
		processusPID, errConv := strconv.Atoi(rgex.FindString(str))
		if errConv != nil {
			log.Panicln(color.RedString("[ERROR]::PID had the wrong format"), errConv)
		}

		aProcess := dash.Processus{
			Name:    processusName,
			PID:     processusPID,
			Version: "N/A",
			Status:  true,
		}

		// Just say here and a bit lower all services we DO want
		switch aProcess.Name {
		case "dockerd":
			docker.Tree = append(docker.Tree, aProcess)
		case "httpd":
			httpd.Tree = append(httpd.Tree, aProcess)
		case "mysqld":
			mysqld.Tree = append(mysqld.Tree, aProcess)
		default: //* To load all the other processes in as "Others"
			// otherProc.Tree = append(otherProc.Tree, aProcess)
		}
		/////////
	}
	// fmt.Println(rgexResult)

	// Explicitly say here and a bit lower all services we DO want
	// If nothing was added we just initialize it with false status
	if len(docker.Tree) == 0 {
		docker.Tree = append(docker.Tree, dash.Processus{
			Name:    "dockerd",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	if len(httpd.Tree) == 0 {
		httpd.Tree = append(httpd.Tree, dash.Processus{
			Name:    "httpd",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	if len(mysqld.Tree) == 0 {
		mysqld.Tree = append(mysqld.Tree, dash.Processus{
			Name:    "mysqld",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	// Final Assignation
	allProcess["Docker"] = docker
	allProcess["Apache"] = httpd
	allProcess["MySQL"] = mysqld
	// allProcess["Others"] = otherProc

	return allProcess
}
