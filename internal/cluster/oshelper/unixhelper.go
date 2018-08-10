package oshelper

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
)

// UnixHelper provied some methods for unix
type UnixHelper struct {
}

func (h UnixHelper) FatalLog(v ...interface{}) {
	log.Fatal(v)
}

func (h UnixHelper) Log(v ...interface{}) {
	log.Print(v)
}

func (h UnixHelper) FileInfo(file string) (string, error) {
	fileInfo, err := os.Stat(file)
	fileInfoJson, _ := json.Marshal(fileInfo)
	return string(fileInfoJson), err
}

func (h UnixHelper) Shell(v ...string) (*exec.Cmd, error) {
	rr := exec.Command(strings.Join(v, " "))
	logString, err := rr.Output()
	if err != nil {
		h.FatalLog(err)
		return rr, err
	}
	h.Log(string(logString))
	return rr, err
}

func (h UnixHelper) GetOs() string {
	return "unix"
}

func (h UnixHelper) CheckDependency(dependecy string) bool {
	_, err := exec.LookPath(dependecy)
	if err != nil {
		h.FatalLog(dependecy, " command not found, kindly check")
		return false
	}

	h.Log("Found", dependecy)
	return true
}
