package config

import (
	"flag"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tenflag "github.com/ten-protocol/go-ten/go/common/flag"
)

func TestCLIFlagTypes(t *testing.T) {
	// Backup the original CommandLine.
	originalFlagSet := flag.CommandLine
	// Create a new FlagSet for testing purposes.
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

	// Defer a function to reset CommandLine after the test.
	defer func() {
		flag.CommandLine = originalFlagSet
	}()

	flags := EnclaveFlags
	err := tenflag.CreateCLIFlags(flags)
	require.NoError(t, err)

	// Set the flags as needed for the test.
	err = flag.CommandLine.Set(LogPathFlag, "string-value")
	require.NoError(t, err)

	err = flag.CommandLine.Set(NodeTypeFlag, "sequencer")
	require.NoError(t, err)

	err = flag.CommandLine.Set(WillAttestFlag, "true")
	require.NoError(t, err)

	err = flag.CommandLine.Set(LogLevelFlag, "123")
	require.NoError(t, err)

	err = flag.CommandLine.Set(MinGasPriceFlag, "3333")
	require.NoError(t, err)

	err = flag.CommandLine.Set(GasBatchExecutionLimit, "222222")
	require.NoError(t, err)

	flag.Parse()

	require.Equal(t, "string-value", flags[LogPathFlag].String())
	require.Equal(t, true, flags[WillAttestFlag].Bool())
	require.Equal(t, 123, flags[LogLevelFlag].Int())
	require.Equal(t, int64(3333), flags[MinGasPriceFlag].Int64())
	require.Equal(t, uint64(222222), flags[GasBatchExecutionLimit].Uint64())

	enclaveConfig, err := newConfig(flags)
	require.NoError(t, err)

	require.Equal(t, "string-value", enclaveConfig.LogPath)
	require.Equal(t, true, enclaveConfig.WillAttest)
	require.Equal(t, 123, enclaveConfig.LogLevel)
	require.Equal(t, big.NewInt(3333), enclaveConfig.MinGasPrice)
	require.Equal(t, uint64(222222), enclaveConfig.GasBatchExecutionLimit)
}

func TestRestrictedMode(t *testing.T) {
	// Backup the original CommandLine.
	originalFlagSet := flag.CommandLine
	// Create a new FlagSet for testing purposes.
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

	// Defer a function to reset CommandLine after the test.
	defer func() {
		flag.CommandLine = originalFlagSet
	}()

	t.Setenv("EDG_TESTMODE", "false")
	t.Setenv("EDG_"+strings.ToUpper(L1ChainIDFlag), "4444")
	t.Setenv("EDG_"+strings.ToUpper(TenChainIDFlag), "1243")
	t.Setenv("EDG_"+strings.ToUpper(TenGenesisFlag), "{}")
	t.Setenv("EDG_"+strings.ToUpper(UseInMemoryDBFlag), "true")
	t.Setenv("EDG_"+strings.ToUpper(ProfilerEnabledFlag), "true")
	t.Setenv("EDG_"+strings.ToUpper(DebugNamespaceEnabledFlag), "true")

	flags := EnclaveFlags
	err := tenflag.CreateCLIFlags(flags)
	require.NoError(t, err)

	err = flag.CommandLine.Set(NodeTypeFlag, "sequencer")
	require.NoError(t, err)

	flag.Parse()

	enclaveConfig, err := NewConfigFromFlags(flags)
	require.NoError(t, err)

	require.Equal(t, int64(4444), enclaveConfig.L1ChainID)
	require.Equal(t, int64(1243), enclaveConfig.TenChainID)
	require.Equal(t, []byte(nil), enclaveConfig.GenesisJSON)
	require.Equal(t, true, enclaveConfig.UseInMemoryDB)
	require.Equal(t, true, enclaveConfig.ProfilerEnabled)
	require.Equal(t, true, enclaveConfig.DebugNamespaceEnabled)
}

func TestRestrictedModeNoCLIDuplication(t *testing.T) {
	// Backup the original CommandLine.
	originalFlagSet := flag.CommandLine
	// Create a new FlagSet for testing purposes.
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

	// Defer a function to reset CommandLine after the test.
	defer func() {
		flag.CommandLine = originalFlagSet
	}()

	t.Setenv("EDG_TESTMODE", "false")
	t.Setenv("EDG_"+strings.ToUpper(L1ChainIDFlag), "4444")
	t.Setenv("EDG_"+strings.ToUpper(TenChainIDFlag), "1243")
	t.Setenv("EDG_"+strings.ToUpper(TenGenesisFlag), "{}")
	t.Setenv("EDG_"+strings.ToUpper(UseInMemoryDBFlag), "true")
	t.Setenv("EDG_"+strings.ToUpper(ProfilerEnabledFlag), "true")
	t.Setenv("EDG_"+strings.ToUpper(DebugNamespaceEnabledFlag), "true")

	flags := EnclaveFlags
	err := tenflag.CreateCLIFlags(flags)
	require.NoError(t, err)

	err = flag.CommandLine.Set(NodeTypeFlag, "sequencer")
	require.NoError(t, err)

	err = flag.CommandLine.Set(L1ChainIDFlag, "5555")
	require.NoError(t, err)

	flag.Parse()

	_, err = NewConfigFromFlags(flags)
	require.Errorf(t, err, "restricted flag was set: l1ChainID")
}
