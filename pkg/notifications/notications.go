package notifications

import "context"

type IEmail interface {
	Send(ctx context.Context, from string, recipient string, sub string, body string) error
}
