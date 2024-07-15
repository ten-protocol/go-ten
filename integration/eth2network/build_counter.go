package eth2network

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const buildNumberFile = "build_number.txt"

func getBuildNumber() (int, error) {
	buildNumber, err := readNextBuildNumber()
	if err != nil {
		return 0, fmt.Errorf("Error getting next build number: %v\n", err)
	}

	err = writeBuildNumber(buildNumber)
	if err != nil {
		return 0, fmt.Errorf("Error saving build number: %v\n", err)
	}
	return buildNumber, nil
}

func readNextBuildNumber() (int, error) {
	data, err := os.ReadFile(buildNumberFile)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil // Start from 1 if file does not exist
		}
		return 0, err
	}

	buildNumber, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	if buildNumber == 99 {
		return 1, err
	}
	return buildNumber + 1, nil
}

func writeBuildNumber(buildNumber int) error {
	return os.WriteFile(buildNumberFile, []byte(fmt.Sprintf("%d", buildNumber)), 0o644)
}
