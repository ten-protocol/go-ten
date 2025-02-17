package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// DebugLogVisibility are the payloads returned when using the DebugGetLogs endpoint
// they are intended for the contract developer to debug visibility issues.
type DebugLogVisibility struct {
	// identify the event
	Address  *gethcommon.Address `json:"contractAddress"`
	EventSig *gethcommon.Hash    `json:"eventSig"`

	// the contract configuration
	AutoContract        bool  `json:"defaultContract"`
	TransparentContract *bool `json:"transparentContract"`

	// the explicit saved configuration for this event type
	EventConfigPublic *bool `json:"eventConfigPublic"`
	Topic1            *bool `json:"topic1Relevant"`
	Topic2            *bool `json:"topic2Relevant"`
	Topic3            *bool `json:"topic3Relevant"`
	Sender            *bool `json:"senderRelevant"`

	// if the event is auto-detected, return which topics will see it
	AutoVisibility *bool `json:"eventAutoVisibility"`
	AutoPublic     *bool `json:"eventAutoPublic"`
	RelAddress1    *bool `json:"topic1AutoRelevant"`
	RelAddress2    *bool `json:"topic2AutoRelevant"`
	RelAddress3    *bool `json:"topic3AutoRelevant"`

	// position of the event
	BlockNumber uint64          `json:"blockNumber" rlp:"-"`
	TxHash      gethcommon.Hash `json:"transactionHash" gencodec:"required" rlp:"-"`
	TxIndex     uint            `json:"transactionIndex" rlp:"-"`
	BlockHash   gethcommon.Hash `json:"blockHash" rlp:"-"`
	Index       uint            `json:"logIndex" rlp:"-"`
}

// LogSubscription is an authenticated subscription to logs.
type LogSubscription struct {
	// ViewingKey - links this subscription request to an externally owed account
	ViewingKey *viewingkey.RPCSignedViewingKey

	// A subscriber-defined filter to apply to the stream of logs.
	Filter *FilterCriteriaJSON
}

func CreateAuthenticatedLogSubscriptionPayload(args []interface{}, vk *viewingkey.ViewingKey) (*LogSubscription, error) {
	logSubscription := &LogSubscription{
		ViewingKey: &viewingkey.RPCSignedViewingKey{
			PublicKey:               vk.PublicKey,
			SignatureWithAccountKey: vk.SignatureWithAccountKey,
			SignatureType:           vk.SignatureType,
		},
	}

	// If there are less than two arguments, it means no filter criteria was passed.
	if len(args) < 2 {
		logSubscription.Filter = &FilterCriteriaJSON{}
		return logSubscription, nil
	}

	filterCriteria, ok := args[1].(FilterCriteria)
	if !ok {
		return nil, fmt.Errorf("invalid subscription")
	}
	fc := SerializableFilterCriteria(filterCriteria)
	logSubscription.Filter = &fc
	return logSubscription, nil
}

// FilterCriteriaJSON is a structure that JSON-serialises to a format that can be successfully deserialised into a
// filters.FilterCriteria object (round-tripping a filters.FilterCriteria to JSON and back doesn't work, due to a
// custom serialiser implemented by filters.FilterCriteria).
type FilterCriteriaJSON struct {
	BlockHash *gethcommon.Hash     `json:"blockHash"`
	FromBlock *rpc.BlockNumber     `json:"fromBlock"`
	ToBlock   *rpc.BlockNumber     `json:"toBlock"`
	Addresses []gethcommon.Address `json:"addresses"`
	Topics    [][]gethcommon.Hash  `json:"topics"`
}

func SerializableFilterCriteria(crit FilterCriteria) FilterCriteriaJSON {
	var from *rpc.BlockNumber
	if crit.FromBlock != nil {
		f := (rpc.BlockNumber)(crit.FromBlock.Int64())
		from = &f
	}

	var to *rpc.BlockNumber
	if crit.ToBlock != nil {
		t := (rpc.BlockNumber)(crit.ToBlock.Int64())
		to = &t
	}

	return FilterCriteriaJSON{
		BlockHash: crit.BlockHash,
		FromBlock: from,
		ToBlock:   to,
		Addresses: crit.Addresses,
		Topics:    crit.Topics,
	}
}

func ToCriteria(jsonCriteria FilterCriteriaJSON) filters.FilterCriteria {
	var from *big.Int
	if jsonCriteria.FromBlock != nil {
		from = big.NewInt(jsonCriteria.FromBlock.Int64())
	}
	var to *big.Int
	if jsonCriteria.ToBlock != nil {
		to = big.NewInt(jsonCriteria.ToBlock.Int64())
	}

	return filters.FilterCriteria{
		BlockHash: jsonCriteria.BlockHash,
		FromBlock: from,
		ToBlock:   to,
		Addresses: jsonCriteria.Addresses,
		Topics:    jsonCriteria.Topics,
	}
}

var errInvalidTopic = errors.New("invalid topic(s)")

// FilterCriteria represents a request to create a new filter.
// Same as ethereum.FilterQuery but with UnmarshalJSON() method.
// duplicated from geth to tweak the unmarshalling
type FilterCriteria ethereum.FilterQuery

// UnmarshalJSON sets *args fields with given data.
func (args *FilterCriteria) UnmarshalJSON(data []byte) error {
	type input struct {
		BlockHash *gethcommon.Hash `json:"blockHash"`
		FromBlock *rpc.BlockNumber `json:"fromBlock"`
		ToBlock   *rpc.BlockNumber `json:"toBlock"`
		Addresses interface{}      `json:"address"`
		Topics    []interface{}    `json:"topics"`
	}

	var raw input
	if err := json.Unmarshal(data, &raw); err != nil {
		// tweak to handle the case when an empty array is passed in by javascript libraries
		if strings.Contains(err.Error(), "cannot unmarshal array") {
			return nil
		}
		return err
	}

	if raw.BlockHash != nil {
		if raw.FromBlock != nil || raw.ToBlock != nil {
			// BlockHash is mutually exclusive with FromBlock/ToBlock criteria
			return errors.New("cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other")
		}
		args.BlockHash = raw.BlockHash
	} else {
		if raw.FromBlock != nil {
			args.FromBlock = big.NewInt(raw.FromBlock.Int64())
		}

		if raw.ToBlock != nil {
			args.ToBlock = big.NewInt(raw.ToBlock.Int64())
		}
	}

	args.Addresses = []gethcommon.Address{}

	if raw.Addresses != nil {
		// raw.Address can contain a single address or an array of addresses
		switch rawAddr := raw.Addresses.(type) {
		case []interface{}:
			for i, addr := range rawAddr {
				if strAddr, ok := addr.(string); ok {
					addr, err := decodeAddress(strAddr)
					if err != nil {
						return fmt.Errorf("invalid address at index %d: %v", i, err)
					}
					args.Addresses = append(args.Addresses, addr)
				} else {
					return fmt.Errorf("non-string address at index %d", i)
				}
			}
		case string:
			addr, err := decodeAddress(rawAddr)
			if err != nil {
				return fmt.Errorf("invalid address: %v", err)
			}
			args.Addresses = []gethcommon.Address{addr}
		default:
			return errors.New("invalid addresses in query")
		}
	}

	// topics is an array consisting of strings and/or arrays of strings.
	// JSON null values are converted to common.Hash{} and ignored by the filter manager.
	if len(raw.Topics) > 0 {
		args.Topics = make([][]gethcommon.Hash, len(raw.Topics))
		for i, t := range raw.Topics {
			switch topic := t.(type) {
			case nil:
				// ignore topic when matching logs

			case string:
				// match specific topic
				top, err := decodeTopic(topic)
				if err != nil {
					return err
				}
				args.Topics[i] = []gethcommon.Hash{top}

			case []interface{}:
				// or case e.g. [null, "topic0", "topic1"]
				for _, rawTopic := range topic {
					if rawTopic == nil {
						// null component, match all
						args.Topics[i] = nil
						break
					}
					if topic, ok := rawTopic.(string); ok {
						parsed, err := decodeTopic(topic)
						if err != nil {
							return err
						}
						args.Topics[i] = append(args.Topics[i], parsed)
					} else {
						return errInvalidTopic
					}
				}
			default:
				return errInvalidTopic
			}
		}
	}

	return nil
}

func decodeAddress(s string) (gethcommon.Address, error) {
	b, err := hexutil.Decode(s)
	if err == nil && len(b) != gethcommon.AddressLength {
		err = fmt.Errorf("hex has invalid length %d after decoding; expected %d for address", len(b), gethcommon.AddressLength)
	}
	return gethcommon.BytesToAddress(b), err
}

func decodeTopic(s string) (gethcommon.Hash, error) {
	b, err := hexutil.Decode(s)
	if err == nil && len(b) != gethcommon.HashLength {
		err = fmt.Errorf("hex has invalid length %d after decoding; expected %d for topic", len(b), gethcommon.HashLength)
	}
	return gethcommon.BytesToHash(b), err
}
