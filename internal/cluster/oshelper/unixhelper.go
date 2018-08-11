package oshelper

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

// UnixHelper provied some methods for unix
type UnixHelper struct {
}

// FatalLog for logging Fatal errors
func (h UnixHelper) FatalLog(v ...interface{}) {
	log.Fatal(v)
}

// Log for logging output
func (h UnixHelper) Log(v ...interface{}) {
	log.Print(v)
}

// FileInfo to get information of a file
func (h UnixHelper) FileInfo(file string) (string, error) {
	fileInfo, err := os.Stat(file)
	fileInfoJSON, _ := json.Marshal(fileInfo)
	return string(fileInfoJSON), err
}

// Shell to execute operstion on the os shell
func (h UnixHelper) Shell(cmd string, v ...string) (*exec.Cmd, error) {
	rr := exec.Command(cmd, v...)
	logString, err := rr.Output()
	if err != nil {
		h.FatalLog(err)
		return rr, err
	}
	h.Log(string(logString))
	return rr, err
}

// CheckDependency checked if the dependency is installed
func (h UnixHelper) CheckDependency(dependecy string) bool {
	_, err := exec.LookPath(dependecy)
	if err != nil {
		h.FatalLog(dependecy, " command not found, kindly check")
		return false
	}

	h.Log("Found", dependecy)
	return true
}
