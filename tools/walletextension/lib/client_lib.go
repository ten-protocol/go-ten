package lib

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	return string(o.userID)
}

func (o *TGLib) Join() error {
	// todo move this to stdlib
	statusCode, userID, err := fasthttp.Get(nil, fmt.Sprintf("%s/v1/join/", o.httpURL))
	if err != nil || statusCode != 200 {
		return fmt.Errorf(fmt.Sprintf("Failed to get userID. Status code: %d, err: %s", statusCode, err))
	}
	o.userID = userID
	return nil
}

func (o *TGLib) RegisterAccount(pk *ecdsa.PrivateKey, addr gethcommon.Address) error {
	// create the registration message
	rawMessageOptions, err := viewingkey.GenerateAuthenticationEIP712RawDataOptions(string(o.userID), integration.TenChainID)
	if err != nil {
		return err
	}
	if len(rawMessageOptions) == 0 {
		return fmt.Errorf("GenerateAuthenticationEIP712RawDataOptions returned 0 options")
	}

	messageHash := crypto.Keccak256(rawMessageOptions[0])
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
		o.httpURL+"/v1/authenticate/?token="+string(o.userID),
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
	personalSignMessage := viewingkey.GeneratePersonalSignMessage(string(o.userID), integration.TenChainID, 1)
	fmt.Println("personalSignMessage: ", personalSignMessage)
	prefixedMessage := fmt.Sprintf(viewingkey.PersonalSignMessagePrefix, len(personalSignMessage), personalSignMessage)
	messageHash := crypto.Keccak256([]byte(prefixedMessage))

	sig, err := crypto.Sign(messageHash, pk)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	sig[64] += 27
	signature := "0x" + hex.EncodeToString(sig)
	fmt.Println("signature: ", signature)
	payload := fmt.Sprintf("{\"signature\": \"%s\", \"address\": \"%s\"}", signature, addr.Hex())

	// issue the registration message
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		o.httpURL+"/v1/authenticate/?token="+string(o.userID),
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
	return fmt.Sprintf("%s/v1/?token=%s", o.httpURL, o.userID)
}

func (o *TGLib) WS() string {
	return fmt.Sprintf("%s/v1/?token=%s", o.wsURL, o.userID)
}
