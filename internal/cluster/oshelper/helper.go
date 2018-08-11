package oshelper

import "os/exec"

// OSHelper to support differnet OS
type OSHelper interface {
	CheckDependency(dependecy string) bool
	FatalLog(...interface{})
	Log(...interface{})
	Shell(cmd string, v ...string) (*exec.Cmd, error)
	FileInfo(file string) (string, error)
}

// NewOSHelper is the Constructor for the OSHelper
func NewOSHelper() OSHelper {
	// TODO check os and return the Helper
	return UnixHelper{}
}
