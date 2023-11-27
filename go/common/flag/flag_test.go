package flag

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringFlagCreation(t *testing.T) {
	expected := "Not the Default"
	flagName := "testString"

	// Create the flag
	flagOutput := String(flagName, "default", "Test string flag")

	// Parse the flag
	os.Args = []string{"cmd", fmt.Sprintf("-%s=%s", flagName, expected)}

	// Parse the flags
	require.NoError(t, Parse())

	// Verify the flag value
	require.Equal(t, expected, *flagOutput)
}

func TestParseInRestrictedMode(t *testing.T) {
	// Set up restricted mode
	t.Setenv("EDG_RESTRICTED", "true")
	defer os.Unsetenv("EDG_RESTRICTED")

	// Create a restricted flag
	flagName := "testFlag"
	flagOutput := RestrictedString(flagName, "default", "Test restricted flag")

	// Mimic setting the environment variable for the restricted flag
	expectedValue := "restrictedValue"
	t.Setenv("EDG_"+strings.ToUpper(flagName), expectedValue)

	defer os.Unsetenv("EDG_" + strings.ToUpper(flagName))

	// Parse the flags
	require.NoError(t, Parse())

	// Verify the flag value
	require.Equal(t, expectedValue, flagOutput.GetString())
}

func TestParseInUnrestrictedMode(t *testing.T) {
	// Ensure unrestricted mode
	os.Unsetenv("EDG_RESTRICTED")

	// Create a regular flag
	flagName := "testUnrestrictedFlag"
	expected := int64(12345)
	// Parse the flag
	os.Args = []string{"cmd", fmt.Sprintf("-%s=%d", flagName, expected)}

	flagOutput := RestrictedInt64(flagName, int64(1), "Test flag")

	// Parse the flags
	require.NoError(t, Parse())

	// Verify the flag value
	require.Equal(t, expected, flagOutput.GetInt64())
}
