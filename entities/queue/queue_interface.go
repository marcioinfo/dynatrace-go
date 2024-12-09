package queue

import "context"

type QueueInterface interface {
	SendMessage(queueUrl string, message string) (err error)
	SendMessageWithContext(ctx context.Context, queueUrl string, message string) (err error)
}
