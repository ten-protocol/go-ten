package eth2network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/codeclysm/extract/v3"
)

const (
	_gethVersion  = "1.16.7"
	_gethHash     = "b9f3a3d9"
	_prysmVersion = "v7.0.0"
	MAC           = "darwin"
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
		// darwin binaries aren't available on geth blobstore so we have to use brew
		if runtime.GOOS == MAC {
			err := installGethViaBrew()
			if err != nil {
				panic(err)
			}
		} else {
			err := checkOrDownloadBinary(gethFileNameVersion, fmt.Sprintf("%s/geth-%s-%s-%s-%s.tar.gz", gethURL, runtime.GOOS, runtime.GOARCH, _gethVersion, _gethHash), true)
			if err != nil {
				panic(err)
			}
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

	client := http.Client{Timeout: 5 * time.Minute}
	// Get the data
	resp, err := client.Get(url) //nolint: gosec, noctx
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

func installGethViaBrew() error {
	expectedDir := path.Join(basepath, _eth2BinariesRelPath, gethFileNameVersion)
	expectedFilePath := path.Join(expectedDir, "geth")

	if fileExists(expectedFilePath) {
		fmt.Printf("Geth already installed at: %s\n", expectedFilePath)
		return nil
	}

	// check if homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("homebrew is not installed. Please install homebrew from https://brew.sh")
	}

	cmd := exec.Command("brew", "install", "ethereum")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install geth via homebrew: %w", err)
	}

	gethPath, err := exec.LookPath("geth")
	if err != nil {
		return fmt.Errorf("geth not found after installation: %w", err)
	}

	// expected directory structure
	if err := os.MkdirAll(expectedDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", expectedDir, err)
	}

	// symlink to the homebrew geth binary
	fmt.Printf("Creating symlink from %s to %s\n", gethPath, expectedFilePath)
	if err := os.Symlink(gethPath, expectedFilePath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Printf("Successfully installed geth via homebrew at: %s\n", expectedFilePath)
	return nil
}
