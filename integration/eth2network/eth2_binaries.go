package eth2network

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	_envUseGethBinary  = "USE_GETH_BINARY"
	_ciBinariesRelPath = "./ci_bin"
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
	//// don't download binaries on each CI
	//if os.Getenv(_envUseGethBinary) == "true" {
	//	return path.Join(basepath, _ciBinariesRelPath), nil
	//}

	// bin folder should exist
	err := os.MkdirAll(path.Join(basepath, _eth2BinariesRelPath), os.ModePerm)
	if err != nil {
		panic(err)
	}

	creationLock.Lock()
	defer creationLock.Unlock()

	goarch := runtime.GOARCH

	// check if the required execs exist, download then otherwise
	if !fileExists(path.Join(basepath, _eth2BinariesRelPath, _gethFileNameVersion)) {
		// geth does not exist
		if runtime.GOARCH == "arm64" {
			goarch = "amd64"
		}

		// tmp download location
		tempFolder, err := os.MkdirTemp("", "")
		if err != nil {
			return "", err
		}

		// download tar file
		downloadFilePath := path.Join(tempFolder, "geth.tar.gz")
		downloadLink := fmt.Sprintf("https://gethstore.blob.core.windows.net/builds/geth-%s-%s-1.10.26-e5eb32ac.tar.gz", runtime.GOOS, goarch)
		err = downloadFile(downloadFilePath, downloadLink)
		if err != nil {
			return "", err
		}

		err = untar(path.Join(basepath, _eth2BinariesRelPath, _gethFileNameVersion), downloadFilePath)
		if err != nil {
			return "", err
		}

		fmt.Printf("Downloaded -  %s\n", _gethFileNameVersion)
	}

	// download prysm files
	for fileName, downloadURL := range map[string]string{
		_prysmBeaconChainFileNameVersion: fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/v3.2.0/beacon-chain-v3.2.0-%s-%s", runtime.GOOS, goarch),
		_prysmCTLFileNameVersion:         fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/v3.2.0/prysmctl-v3.2.0-%s-%s", runtime.GOOS, goarch),
		_prysmValidatorFileNameVersion:   fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/v3.2.0/validator-v3.2.0-%s-%s", runtime.GOOS, goarch),
	} {
		expectedFilePath := path.Join(basepath, _eth2BinariesRelPath, fileName)
		if !fileExists(expectedFilePath) {
			// download tar file
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

func untar(dst string, path string) error {
	openFile, err := os.Open(path)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(openFile)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case errors.Is(err, io.EOF):
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// ignore non files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// it's a file create it
		fileName := strings.Split(header.Name, "/")[1]
		if fileName == "geth" {
			f, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil { //nolint: gosec
				return err
			}

			f.Close()
		}
	}
}
