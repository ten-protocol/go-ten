package eth2network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/codeclysm/extract/v3"
)

const (
	_gethVersion  = "1.14.6"
	_gethHash     = "aadddf3a"
	_prysmVersion = "v5.0.4"
	//_gethVersion  = "1.15.5"
	//_gethHash     = "4263936a"
	//_prysmVersion = "v5.3.0"
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

	var wg sync.WaitGroup
	prysmaticURL := fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/%s/", _prysmVersion)
	gethURL := "https://gethstore.blob.core.windows.net/builds"
	wg.Add(4)
	go func() {
		defer wg.Done()
		err := checkOrDownloadBinary(prysmBeaconChainFileNameVersion, fmt.Sprintf("%s%s", prysmaticURL, prysmBeaconChainFileNameVersion), false)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := checkOrDownloadBinary(prysmCTLFileNameVersion, fmt.Sprintf("%s%s", prysmaticURL, prysmCTLFileNameVersion), false)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := checkOrDownloadBinary(prysmValidatorFileNameVersion, fmt.Sprintf("%s%s", prysmaticURL, prysmValidatorFileNameVersion), false)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := checkOrDownloadBinary(gethFileNameVersion, fmt.Sprintf("%s/geth-%s-%s-%s-%s.tar.gz", gethURL, runtime.GOOS, runtime.GOARCH, _gethVersion, _gethHash), true)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	return path.Join(basepath, _eth2BinariesRelPath), nil
}

func checkOrDownloadBinary(fileName string, url string, unTar bool) error {
	expectedFilePath := path.Join(basepath, _eth2BinariesRelPath, fileName)
	if fileExists(expectedFilePath) {
		return nil
	}
	suffix := ""
	if unTar {
		suffix = ".tar.gz"
	}
	err := downloadFile(expectedFilePath+suffix, url)
	if err != nil {
		return err
	}
	if unTar {
		f, err := os.Open(expectedFilePath + suffix)
		if err != nil {
			return err
		}
		err = extract.Gz(context.TODO(), f, path.Join(basepath, _eth2BinariesRelPath), nil)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Downloaded - %s\n", expectedFilePath)
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
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
		os.Remove(filepath)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
