package eth2network

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	_gethVersion  = "1.10.26"
	_prysmVersion = "v3.2.0"
)

var (
	// prevents issues when calling from different packages/directories
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)

	creationLock sync.Mutex // makes sure there isn't two creations running at the same time
)

// EnsureBinariesExist makes sure node binaries exist, returns the base path where binaries exist
// Downloads any missing binaries
func EnsureBinariesExist() (string, error) {
	creationLock.Lock()
	defer creationLock.Unlock()

	// bin folder should exist
	err := os.MkdirAll(path.Join(basepath, _eth2BinariesRelPath), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// build geth
	if !fileExists(path.Join(basepath, _eth2BinariesRelPath, _gethFileNameVersion)) {
		gethScript := path.Join(basepath, "./build_geth_binary.sh")
		cmd := exec.Command("bash", gethScript, fmt.Sprintf("%s=%s", "--version", "v"+_gethVersion))
		cmd.Stderr = os.Stderr

		if out, err := cmd.Output(); err != nil {
			fmt.Printf("%s\n", out)
			return "", err
		}
	}

	// download prysm files
	for fileName, downloadURL := range map[string]string{
		_prysmBeaconChainFileNameVersion: fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-%s-%s", _prysmVersion, _prysmVersion, runtime.GOOS, runtime.GOARCH),
		_prysmCTLFileNameVersion:         fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/%s/prysmctl-%s-%s-%s", _prysmVersion, _prysmVersion, runtime.GOOS, runtime.GOARCH),
		_prysmValidatorFileNameVersion:   fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-%s-%s", _prysmVersion, _prysmVersion, runtime.GOOS, runtime.GOARCH),
	} {
		expectedFilePath := path.Join(basepath, _eth2BinariesRelPath, fileName)
		if !fileExists(expectedFilePath) {
			err := downloadFile(expectedFilePath, downloadURL)
			if err != nil {
				return "", err
			}
			fmt.Printf("Downloaded - %s\n", fileName)
		}
	}

	return path.Join(basepath, _eth2BinariesRelPath), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadFile(filepath string, url string) error {
	fmt.Printf("Downloading: %s\n", url)
	// Create the file
	out, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o777)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url) //nolint: gosec, noctx
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
