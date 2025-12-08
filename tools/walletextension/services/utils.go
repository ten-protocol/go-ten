package services

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// ReturnDefaultUserAndAccount creates a new in-memory user and a corresponding account.
// Nothing is persisted to storage. Useful for anonymous/public flows.
func ReturnDefaultUserAndAccount(config *common.Config) (*common.GWUser, error) {
	// generate a fresh viewing key
	defaultUserVK := "5b7db1a436d96273b4ebb8a5bb28d59f28d1d54810b723dd6e03731ec335d10c" // hardcoded viewing key for the default user - remove after proper public access is implemented
	defaultUserVKBytes, err := hex.DecodeString(defaultUserVK)
	if err != nil {
		return nil, fmt.Errorf("failed to decode default VK hex: %w", err)
	}
	vk, err := crypto.ToECDSA(defaultUserVKBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert default VK bytes to private key: %w", err)
	}
	vkEcies := ecies.ImportECDSA(vk)

	// derive user ID from the viewing key
	userID := viewingkey.CalculateUserID(common.PrivateKeyToCompressedPubKey(vkEcies))

	// build an in-memory GWUser (no persistence)
	user := &common.GWUser{
		ID:          userID,
		Accounts:    make(map[gethcommon.Address]*common.GWAccount),
		UserKey:     crypto.FromECDSA(vkEcies.ExportECDSA()),
		SessionKeys: make(map[gethcommon.Address]*common.GWSessionKey),
	}

	userAddress := crypto.PubkeyToAddress(vk.PublicKey)
	msg, err := viewingkey.GenerateMessage(user.ID, int64(config.TenChainID), 1, viewingkey.EIP712Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message. Cause %w", err)
	}

	msgHash, err := viewingkey.GetMessageHash(msg, viewingkey.EIP712Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message hash. Cause %w", err)
	}

	// current signature is valid - return account address
	sig, err := crypto.Sign(msgHash, vk)
	if err != nil {
		return nil, fmt.Errorf("cannot sign message with session key. Cause %w", err)
	}

	// create an account that signs over the userID
	account := &common.GWAccount{
		User:          user,
		Address:       &userAddress,
		Signature:     sig,
		SignatureType: viewingkey.EIP712Signature,
	}

	user.Accounts[userAddress] = account

	return user, nil
}

func SafeArgsForLogging(args []any) string {
	if len(args) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, arg := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		safeStringify(&b, reflect.ValueOf(arg), 0)
	}
	b.WriteByte(']')
	return b.String()
}

const maxDepth = 5

func safeStringify(b *strings.Builder, rv reflect.Value, depth int) {
	if depth > maxDepth {
		b.WriteString("...")
		return
	}

	if !rv.IsValid() {
		b.WriteString("<nil>")
		return
	}

	switch rv.Kind() {
	case reflect.Interface:
		if rv.IsZero() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		safeStringify(b, rv.Elem(), depth+1)
		return

	case reflect.Ptr:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		safeStringify(b, rv.Elem(), depth+1)
		return

	case reflect.Slice:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		b.WriteByte('[')
		for i := 0; i < rv.Len() && i < 10; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			safeStringify(b, rv.Index(i), depth+1)
		}
		if rv.Len() > 10 {
			fmt.Fprintf(b, "...+%d more", rv.Len()-10)
		}
		b.WriteByte(']')
		return

	case reflect.Map:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		b.WriteString("map[")
		iter := rv.MapRange()
		count := 0
		for iter.Next() && count < 5 {
			if count > 0 {
				b.WriteString(", ")
			}
			safeStringify(b, iter.Key(), depth+1)
			b.WriteByte(':')
			safeStringify(b, iter.Value(), depth+1)
			count++
		}
		b.WriteByte(']')
		return

	case reflect.Struct:
		b.WriteByte('{')
		t := rv.Type()
		for i := 0; i < rv.NumField() && i < 10; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(t.Field(i).Name)
			b.WriteByte(':')
			field := rv.Field(i)
			if field.CanInterface() {
				safeStringify(b, field, depth+1)
			} else {
				b.WriteString("<unexported>")
			}
		}
		b.WriteByte('}')
		return

	case reflect.Chan, reflect.Func:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
		} else {
			fmt.Fprintf(b, "<%s>", rv.Type())
		}
		return

	case reflect.String:
		fmt.Fprintf(b, "%q", rv.String())
	case reflect.Bool:
		fmt.Fprintf(b, "%t", rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(b, "%d", rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(b, "%d", rv.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(b, "%g", rv.Float())
	default:
		fmt.Fprintf(b, "<%s>", rv.Type())
	}
}
