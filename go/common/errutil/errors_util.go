package errutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	// ErrNotFound must equal Geth's not-found error. This is because some Geth components we use throw the latter, and
	// we want to be able to catch both types in a single error-check.
	ErrNotFound      = ethereum.NotFound
	ErrAlreadyExists = errors.New("already exists")
	ErrNoImpl        = errors.New("not implemented")

	// Standard errors that can be returned from block submission

	ErrBlockAlreadyProcessed       = errors.New("block already processed")
	ErrBlockAncestorNotFound       = errors.New("block ancestor not found")
	ErrBlockForBatchNotFound       = errors.New("block for batch not found")
	ErrAncestorBatchNotFound       = errors.New("parent for batch not found")
	ErrCrossChainBundleRepublished = errors.New("root already added to the message bus")
	ErrCrossChainBundleNoBatches   = errors.New("no batches for cross chain bundle")
	ErrNoNextRollup                = errors.New("no next rollup")
	ErrRollupForkMismatch          = errors.New("rollup fork mismatch")
	ErrNoBundleToPublish           = errors.New("no bundle to publish")
	ErrBlockBindingMismatch        = errors.New("Block binding mismatch")
)

// BlockRejectError is used as a standard format for error response from enclave for block submission errors
// The L1 Head hash tells the host what the enclave knows as the canonical chain head, so it can feed it the appropriate block.
type BlockRejectError struct {
	L1Head  gethcommon.Hash
	Wrapped error
}

func (r BlockRejectError) Error() string {
	head := "N/A"
	if r.L1Head != (gethcommon.Hash{}) {
		head = r.L1Head.String()
	}
	return fmt.Sprintf("%s l1Head=%s", r.Wrapped.Error(), head)
}

func (r BlockRejectError) Unwrap() error {
	return r.Wrapped
}

// Is implementation supports the errors.Is() behaviour. The error being sent over the wire means we lose the typing
// on the `Wrapped` error, but comparing the string helps in our use cases of standard exported errors above
func (r BlockRejectError) Is(err error) bool {
	return strings.Contains(r.Error(), err.Error()) || errors.Is(err, r.Wrapped)
}
