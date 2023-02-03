package traffic

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

// SuccessesVerifier checks that all actions carried out were successful
func SuccessesVerifier(data RunData, _ networktest.NetworkConnector) error {
	failCount := 0
	for _, action := range data.ActionEvents() {
		if !action.Success() {
			failCount++
		}
	}
	successesDesc := fmt.Sprintf("%d actions executed (%d failed)", len(data.ActionEvents()), failCount)
	fmt.Println(successesDesc)
	testlog.Logger().Info(successesDesc)
	if failCount > 0 {
		return fmt.Errorf("FAIL %s - %s", data.RunDescription(), successesDesc)
	}

	return nil
}
