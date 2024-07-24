package eth2network

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

// GetBuildNumber retrieves and increments the build number from an environment variable.
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

	// Increment the build number until an unused folder is found
	for {
		buildPath := path.Join(basepath, "../.build/eth2", strconv.Itoa(buildNumber))
		if _, err := os.Stat(buildPath); os.IsNotExist(err) {
			break
		}
		buildNumber++
		if buildNumber > 9999 {
			buildNumber = 1
		}
	}

	err = os.Setenv("BUILD_NUMBER", strconv.Itoa(buildNumber))
	if err != nil {
		return 0, fmt.Errorf("Error setting build number: %v\n", err)
	}

	return buildNumber, nil
}
