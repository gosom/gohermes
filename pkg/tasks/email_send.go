package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/notifications"
	"github.com/hibiken/asynq"
)

const TypeEmailDelivery = "email:delivery"

type EmailDeliveryPayload struct {
	From      string
	Recipient string
	Sub       string
	Body      string
}

func NewEmailDeliveryTask(from, recipient, sub, body string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{
		Recipient: recipient, Sub: sub, Body: body,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

type EmailDeliveryProcessor struct {
	di      *container.ServiceContainer
	backend notifications.IEmail
}

func NewEmailDeliveryProcessor(di *container.ServiceContainer) (*EmailDeliveryProcessor, error) {
	ans := EmailDeliveryProcessor{}
	ans.di = di
	ibackend, err := di.GetService("email_backend")
	if err != nil {
		return nil, fmt.Errorf("cannot find email_backend: %w", err)
	}
	ans.backend = ibackend.(notifications.IEmail)
	return &ans, nil
}

func (o *EmailDeliveryProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	err := o.backend.Send(ctx, p.From, p.Recipient, p.Sub, p.Body)
	if err != nil {
		o.di.Logger.Error().Msgf("faild to sent email to %s: %s", p.Recipient, err.Error())
		return err
	}
	o.di.Logger.Info().Msgf("sent email to %s", p.Recipient)
	return nil
}
