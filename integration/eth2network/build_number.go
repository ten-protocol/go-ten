package eth2network

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

// Need this due to character limit on unix path when deploying the geth network on azure and integration tests. We
// incremenet the build number up to 9999 and then reset to 1, overwriting any existing files.
func getBuildNumber() (int, error) {
	buildNumberStr := os.Getenv("BUILD_NUMBER")
	var buildNumber int
	var err error

	if buildNumberStr == "" {
		buildNumber = 1
	} else {
		buildNumber, err = strconv.Atoi(buildNumberStr)
		if err != nil {
			return 0, fmt.Errorf("Error converting build number: %v\n", err)
		}
	}

	// increment the build number until an unused folder is found
	for {
		buildPath := path.Join(basepath, "../.build/eth2", strconv.Itoa(buildNumber))
		if _, err := os.Stat(buildPath); os.IsNotExist(err) {
			break
		}
		buildNumber++
		if buildNumber > 99 {
			return 0, fmt.Errorf("Error: no available build number from 1-99. Delete some build folders!\n")
		}
	}

	err = os.Setenv("BUILD_NUMBER", strconv.Itoa(buildNumber))
	if err != nil {
		return 0, fmt.Errorf("Error setting build number: %v\n", err)
	}

	return buildNumber, nil
}
