package dash

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kataras/iris"
)

// // Process

// Process is a process we can start or stop on the server
type Process interface {
	Init() error
	Start() error
	Stop() error
	GetPID() int
	GetVersion() string
	GetStatus() bool
}

// Apache is an apache server, apache2 or httpd
type Apache struct {
	Name    string `json:"service"`
	Version string `json:"version"`
	Status  bool   `json:"status"`
}

// Processus is just a simple processus
type Processus struct {
	Name    string `json:"name"`
	PID     int    `json:"pid"`
	Version string `json:"version"`
	Status  bool   `json:"status"`
}

// TreeProcess is the "arborescence" all processussesssss
type TreeProcess struct {
	Tree []Processus `json:"tree"`
}

// Init initializes the processus from its name
func (p *Processus) Init() error {
	//TODO: Make first the func that fetch all process
	//TODO Then look for it in the slice returned by the fetch
	//TODO Init to 0 value if it's not there
	return nil
}

// Start the apache service
func (p *Processus) Start() error {
	return nil
}

// Stop is the action to stop the apache service
//! Do not kill process, just stop service bitch
//* Accesses by PUT method
func (p *Processus) Stop() error {
	// new command
	// cmd := exec.Command("kill", p.PID)
	q := fmt.Sprintf("kill %s", p.GetPID())
	cmd := exec.Command("sudo", "bash", q)
	err := cmd.Run()
	if err != nil {
		log.Panicf("[ERROR]::Could not kill the process PID[%d]\n\t%v\n", p.PID, err)
	}

	return nil
}

// GetPID get the pid of the apache service
func (p *Processus) GetPID() string {
	return strconv.Itoa(p.PID)
}

// GetVersion return the version of the apache service
func (p *Processus) GetVersion() string {
	return p.Version
}

// GetStatus return false if the process is down, true otherwise
func (p *Processus) GetStatus() bool {
	return false
}

// // // Service
// type Service interface {
// 	Get() []Process
// }
//
// // ServerService is a struct
// type ServerService struct {
// 	items map[string]
// }

//= NON ATTACHED METHODS

// LoadProcess loads the processes
func LoadProcess() iris.Map {
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
	httpd := TreeProcess{
		Tree: make([]Processus, 0),
	}
	dockerd := TreeProcess{
		Tree: make([]Processus, 0),
	}
	mysqld := TreeProcess{
		Tree: make([]Processus, 0),
	}
	mongod := TreeProcess{
		Tree: make([]Processus, 0),
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

		aProcess := Processus{
			Name:    processusName,
			PID:     processusPID,
			Version: "N/A",
			Status:  true,
		}

		//* Just say here and a bit lower all services we DO want
		switch aProcess.Name {
		case "dockerd":
			dockerd.Tree = append(dockerd.Tree, aProcess)
		case "httpd":
			httpd.Tree = append(httpd.Tree, aProcess)
		case "mysqld":
			mysqld.Tree = append(mysqld.Tree, aProcess)
		case "mongod":
			mongod.Tree = append(mongod.Tree, aProcess)
		default: //* To load all the other processes in as "Others"
			// otherProc.Tree = append(otherProc.Tree, aProcess)
		}
		/////////
	}
	// fmt.Println(rgexResult)

	// Explicitly say here and a bit lower all services we DO want
	// If nothing was added we just initialize it with false status
	if len(dockerd.Tree) == 0 {
		dockerd.Tree = append(dockerd.Tree, Processus{
			Name:    "dockerd",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	if len(httpd.Tree) == 0 {
		httpd.Tree = append(httpd.Tree, Processus{
			Name:    "httpd",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	if len(mysqld.Tree) == 0 {
		mysqld.Tree = append(mysqld.Tree, Processus{
			Name:    "mysqld",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	if len(mongod.Tree) == 0 {
		mongod.Tree = append(mongod.Tree, Processus{
			Name:    "mongod",
			PID:     0,
			Version: "N/A",
			Status:  false,
		})
	}
	// Final Assignation
	allProcess["Docker"] = dockerd
	allProcess["Apache"] = httpd
	allProcess["MySQL"] = mysqld
	allProcess["MongoDB"] = mongod
	// allProcess["Others"] = otherProc

	return allProcess
}
