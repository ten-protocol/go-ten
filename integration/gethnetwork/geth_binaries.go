package gethnetwork

import (
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	shCmd             = "sh"
	gethScriptPathRel = "./build_geth_binary.sh"
	gethBinaryPathRel = "../.build/geth_bin/geth"
	versionFlag       = "--version"

	LatestVersion = "v1.10.17" // geths release version
)

var (
	// prevents issues when calling from different packages/directories
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)

	creationLock sync.Mutex // makes sure there isn't two creations running at the same time
)

// EnsureBinariesExist makes sure geth binary exist, returns path where geth binaries exist
// Runs the build_geth_binary.sh while handling source directory call
func EnsureBinariesExist(version string) (string, error) {
	creationLock.Lock()
	defer creationLock.Unlock()

	gethScript := path.Join(basepath, gethScriptPathRel)
	_, err := exec.Command(shCmd, gethScript, fmt.Sprintf("%s=%s", versionFlag, version)).Output()
	if err != nil {
		return "", err
	}
	return path.Join(basepath, fmt.Sprintf("%s-%s", gethBinaryPathRel, version)), nil
}
