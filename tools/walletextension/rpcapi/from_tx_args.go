package rpcapi

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
)

func searchFromAndData(possibleAddresses []*common.Address, args gethapi.TransactionArgs) *common.Address {
	if args.From != nil && (*args.From != common.Address{}) {
		return args.From
	}

	if args.Data == nil {
		return nil
	}

	// since the "from" field is not mandatory, we try to find a matching address in the data field
	addressesMap := toMap(possibleAddresses)
	return searchDataFieldForAccount(addressesMap, *args.Data)
}

func searchDataFieldForAccount(addressesMap map[common.Address]*common.Address, data []byte) *common.Address {
	hexEncodedData := hexutils.BytesToHex(data)

	// We check that the data field is long enough before removing the leading "0x" (1 bytes/2 chars) and the method ID
	// (4 bytes/8 chars).
	if len(hexEncodedData) < 10 {
		return nil
	}
	hexEncodedData = hexEncodedData[10:]

	// We split up the arguments in the `data` field.
	var dataArgs []string
	for i := 0; i < len(hexEncodedData); i += ethCallPaddedArgLen {
		if i+ethCallPaddedArgLen > len(hexEncodedData) {
			break
		}
		dataArgs = append(dataArgs, hexEncodedData[i:i+ethCallPaddedArgLen])
	}

	// We iterate over the arguments, looking for an argument that matches a viewing key address
	for _, dataArg := range dataArgs {
		// If the argument doesn't have the correct padding, it's not an address.
		if !strings.HasPrefix(dataArg, ethCallAddrPadding) {
			continue
		}

		maybeAddress := common.HexToAddress(dataArg[len(ethCallAddrPadding):])
		if _, ok := addressesMap[maybeAddress]; ok {
			return &maybeAddress
		}
	}

	return nil
}

func toMap(possibleAddresses []*common.Address) map[common.Address]*common.Address {
	addresses := map[common.Address]*common.Address{}
	for i := range possibleAddresses {
		addresses[*possibleAddresses[i]] = possibleAddresses[i]
	}
	return addresses
}
