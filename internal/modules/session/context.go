package session

import (
	"context"
	"fmt"
)

type AccountIDContextKey struct{}

// func WithAccountID(ctx context.Context, accountID string) context.Context {
// 	return context.WithValue(ctx, AccountIDContextKey{}, accountID)
// }

func (service *Service) GetAccountID(ctx context.Context) (string, bool) {
	accountID := ctx.Value(AccountIDContextKey{})

	fmt.Println("AccountID", accountID)

	return accountID.(string), accountID != nil
}
