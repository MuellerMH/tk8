package oshelper

import "os/exec"

type OSHelper interface {
	GetOs() string
	CheckDependency(dependecy string) bool
	FatalLog(...interface{})
	Log(...interface{})
	Shell(...string) (*exec.Cmd, error)
	FileInfo(file string) (string, error)
}

func NewOSHelper() OSHelper {
	// TODO check os and return the Helper
	return UnixHelper{}
}
