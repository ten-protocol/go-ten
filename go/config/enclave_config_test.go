package config

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	tenflag "github.com/ten-protocol/go-ten/go/common/flag"
)

func TestFromFlags(t *testing.T) {
	// Backup the original CommandLine.
	originalFlagSet := flag.CommandLine
	// Create a new FlagSet for testing purposes.
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

	// Defer a function to reset CommandLine after the test.
	defer func() {
		flag.CommandLine = originalFlagSet
	}()

	flags := EnclaveFlags()
	err := tenflag.CreateCLIFlags(flags)
	require.NoError(t, err)

	// Set the flags as needed for the test.
	err = flag.CommandLine.Set(HostIDFlag, "string-value")
	require.NoError(t, err)

	err = flag.CommandLine.Set(NodeTypeFlag, "sequencer")
	require.NoError(t, err)

	err = flag.CommandLine.Set(WillAttestFlag, "true")
	require.NoError(t, err)

	err = flag.CommandLine.Set(LogLevelFlag, "123")
	require.NoError(t, err)

	err = flag.CommandLine.Set(MinGasPriceFlag, "3333")
	require.NoError(t, err)

	err = flag.CommandLine.Set(L2GasLimitFlag, "222222")
	require.NoError(t, err)

	flag.Parse()

	require.Equal(t, "string-value", flags[HostIDFlag].String())
	require.Equal(t, true, flags[WillAttestFlag].Bool())
	require.Equal(t, 123, flags[LogLevelFlag].Int())
	require.Equal(t, int64(3333), flags[MinGasPriceFlag].Int64())
	require.Equal(t, uint64(222222), flags[L2GasLimitFlag].Uint64())

	_, err = newConfig(flags)
	require.NoError(t, err)
}
