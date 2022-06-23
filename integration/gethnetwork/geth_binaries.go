package gethnetwork

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	shCmd                        = "bash"
	gethScriptPathRel            = "./build_geth_binary.sh"
	gethBinaryPathRel            = "../.build/geth_bin/geth"
	gethPrecompiledBinaryPathRel = "./geth_bin/geth-v1.10.17"
	versionFlag                  = "--version"
	envUseGethBinary             = "USE_GETH_BINARY"

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
	if os.Getenv(envUseGethBinary) == "true" {
		return path.Join(basepath, gethPrecompiledBinaryPathRel), nil
	}

	creationLock.Lock()
	defer creationLock.Unlock()

	gethScript := path.Join(basepath, gethScriptPathRel)
	cmd := exec.Command(shCmd, gethScript, fmt.Sprintf("%s=%s", versionFlag, version))
	cmd.Stderr = os.Stderr

	if out, err := cmd.Output(); err != nil {
		fmt.Printf("%s\n", out)
		return "", err
	}
	return path.Join(basepath, fmt.Sprintf("%s-%s", gethBinaryPathRel, version)), nil
}
