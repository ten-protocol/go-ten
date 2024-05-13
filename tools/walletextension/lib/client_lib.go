package lib

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/integration"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/valyala/fasthttp"
)

type TGLib struct {
	httpURL string
	wsURL   string
	userID  []byte
}

func NewTenGatewayLibrary(httpURL, wsURL string) *TGLib {
	return &TGLib{
		httpURL: httpURL,
		wsURL:   wsURL,
	}
}

func (o *TGLib) UserID() string {
	return hexutils.BytesToHex(o.userID)
}

func (o *TGLib) UserIDBytes() []byte {
	return o.userID
}

func (o *TGLib) Join() error {
	// todo move this to stdlib
	statusCode, userID, err := fasthttp.Get(nil, fmt.Sprintf("%s/v1/join/", o.httpURL))
	if err != nil || statusCode != 200 {
		return fmt.Errorf(fmt.Sprintf("Failed to get userID. Status code: %d, err: %s", statusCode, err))
	}
	o.userID = hexutils.HexToBytes(string(userID))
	return nil
}

func (o *TGLib) RegisterAccount(pk *ecdsa.PrivateKey, addr gethcommon.Address) error {
	// create the registration message
	message, err := viewingkey.GenerateMessage(o.userID, integration.TenChainID, 1, viewingkey.EIP712Signature)
	if err != nil {
		return err
	}

	messageHash, err := viewingkey.GetMessageHash(message, viewingkey.EIP712Signature)
	if err != nil {
		return fmt.Errorf("failed to get message hash: %w", err)
	}

	sig, err := crypto.Sign(messageHash, pk)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	sig[64] += 27
	signature := "0x" + hex.EncodeToString(sig)
	payload := fmt.Sprintf("{\"signature\": \"%s\", \"address\": \"%s\"}", signature, addr.Hex())

	// issue the registration message
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		o.httpURL+"/v1/authenticate/?token=0x"+hexutils.BytesToHex(o.userID),
		strings.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("unable to create request - %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to issue request - %w", err)
	}

	defer response.Body.Close()
	r, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read response - %w", err)
	}
	if string(r) != "success" {
		return fmt.Errorf("expected success, got %s", string(r))
	}
	return nil
}

func (o *TGLib) RegisterAccountPersonalSign(pk *ecdsa.PrivateKey, addr gethcommon.Address) error {
	// create the registration message
	message, err := viewingkey.GenerateMessage(o.userID, integration.TenChainID, viewingkey.PersonalSignVersion, viewingkey.PersonalSign)
	if err != nil {
		return err
	}

	messageHash, err := viewingkey.GetMessageHash(message, viewingkey.PersonalSign)
	if err != nil {
		return fmt.Errorf("failed to get message hash: %w", err)
	}

	sig, err := crypto.Sign(messageHash, pk)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	sig[64] += 27
	signature := "0x" + hex.EncodeToString(sig)
	payload := fmt.Sprintf("{\"signature\": \"%s\", \"address\": \"%s\", \"type\": \"%s\"}", signature, addr.Hex(), "Personal")

	// issue the registration message
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		o.httpURL+"/v1/authenticate/?token=0x"+hexutils.BytesToHex(o.userID),
		strings.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("unable to create request - %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to issue request - %w", err)
	}

	defer response.Body.Close()
	r, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read response - %w", err)
	}
	if string(r) != "success" {
		return fmt.Errorf("expected success, got %s", string(r))
	}
	return nil
}

func (o *TGLib) HTTP() string {
	return fmt.Sprintf("%s/v1/?token=%s", o.httpURL, hexutils.BytesToHex(o.userID))
}

func (o *TGLib) WS() string {
	return fmt.Sprintf("%s/v1/?token=%s", o.wsURL, hexutils.BytesToHex(o.userID))
}
