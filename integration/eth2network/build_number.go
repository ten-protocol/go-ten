package eth2network

import (
	"fmt"
	"os"
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
		if buildNumber == 99 {
			buildNumber = 1
		} else {
			buildNumber++
		}
	}

	err = os.Setenv("BUILD_NUMBER", strconv.Itoa(buildNumber))
	if err != nil {
		return 0, fmt.Errorf("Error setting build number: %v\n", err)
	}

	return buildNumber, nil
}
