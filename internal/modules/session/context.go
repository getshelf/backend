package session

import (
	"context"
)

type AccountIDContextKey struct{}

func (service *Service) GetAccountID(ctx context.Context) (string, bool) {
	accountID, ok := ctx.Value(AccountIDContextKey{}).(string)

	return accountID, ok
}
