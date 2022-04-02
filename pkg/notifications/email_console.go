package notifications

import (
	"context"
	"fmt"
)

type EmailBackendConsole struct {
}

func (o *EmailBackendConsole) Send(ctx context.Context, from, recipient string, sub string, body string) error {
	printLnPadded(fmt.Sprintf("Email from: %s to: %s", from, recipient))
	fmt.Println("")
	fmt.Printf("Subject: %s\n", sub)
	fmt.Println("")
	fmt.Printf("Body:")
	fmt.Println(body)
	fmt.Println("")
	printLnPadded("Email finish")
	return nil
}

func printLnPadded(s string) {
	w := 110
	fmt.Printf(fmt.Sprintf("%%-%ds\n", w/2), fmt.Sprintf(fmt.Sprintf("%%%ds", w/2), s))
}
