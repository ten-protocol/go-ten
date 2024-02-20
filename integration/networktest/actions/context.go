package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"
)

// KeyNumberOfTestUsers key to an int representing number of test users created/available
// (useful for test actions that want to run for all users without having to duplicate config)
var KeyNumberOfTestUsers = ActionKey("numberOfTestUsers")

// ActionKey is the type for all test data stored in the context. Go documentation recommends using a typed key rather than string to avoid conflicts.
type ActionKey string

func storeTestUser(ctx context.Context, userNumber int, user userwallet.User) context.Context {
	return context.WithValue(ctx, userKey(userNumber), user)
}

func FetchTestUser(ctx context.Context, userNumber int) (userwallet.User, error) {
	u := ctx.Value(userKey(userNumber))
	if u == nil {
		return nil, fmt.Errorf("no userWallet found in context for userNumber=%d", userNumber)
	}
	user, ok := u.(userwallet.User)
	if !ok {
		return nil, fmt.Errorf("user retrieved from context was not of expected type userWallet for userNumber=%d type=%T", userNumber, u)
	}
	return user, nil
}

func userKey(number int) ActionKey {
	return ActionKey(fmt.Sprintf("testUser%d", number))
}

func FetchNumberOfTestUsers(ctx context.Context) (int, error) {
	n := ctx.Value(KeyNumberOfTestUsers)
	if n == nil {
		return 0, errors.New("expected at least one test user to be setup but number of test users was not set")
	}
	numUsers, ok := n.(int)
	if !ok || numUsers == 0 {
		return 0, fmt.Errorf("expected number of users context value to be non-zero int but it was (%T) %v", n, n)
	}
	return numUsers, nil
}

func FetchBigInt(ctx context.Context, key ActionKey) (*big.Int, error) {
	v := ctx.Value(key)
	if v == nil {
		return nil, fmt.Errorf("no value found for key=%s", key)
	}
	typedVal, ok := v.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("expected value for key=%s to be *big.Int but was (%T) %v", key, v, v)
	}
	return typedVal, nil
}
