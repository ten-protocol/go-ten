package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"
)

// int representing number of test users created/available (useful for tests that want to run for all users)
var keyNumberOfTestUsers = ActionKey("numberOfTestUsers")

// ActionKey is the type for all test data stored in the context. Go documentation recommends using a typed key rather than string to avoid conflicts.
type ActionKey string

func storeTestUser(ctx context.Context, userNumber int, user *userwallet.UserWallet) context.Context {
	return context.WithValue(ctx, userKey(userNumber), user)
}

func FetchTestUser(ctx context.Context, userNumber int) (*userwallet.UserWallet, error) {
	u := ctx.Value(userKey(userNumber))
	if u == nil {
		return nil, fmt.Errorf("no UserWallet found in context for userNumber=%d", userNumber)
	}
	user, ok := u.(*userwallet.UserWallet)
	if !ok {
		return nil, fmt.Errorf("user retrieved from context was not of expected type UserWallet for userNumber=%d type=%T", userNumber, u)
	}
	return user, nil
}

func userKey(number int) ActionKey {
	return ActionKey(fmt.Sprintf("testUser%d", number))
}

func FetchNumberOfTestUsers(ctx context.Context) (int, error) {
	n := ctx.Value(keyNumberOfTestUsers)
	if n == nil {
		return 0, errors.New("expected at least one test user to be setup but number of test users was not set")
	}
	numUsers, ok := n.(int)
	if !ok || numUsers == 0 {
		return 0, fmt.Errorf("expected number of users context value to be non-zero int but it was (%T) %v", n, n)
	}
	return numUsers, nil
}
