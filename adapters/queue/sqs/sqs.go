package sqs

import (
	"context"
	"os"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"

	aws_adapter "payment-layer-card-api/adapters/aws"
)

type SqsAWS struct {
	awsConnection *aws_adapter.AWSConnection
	groupId       string
}

func NewSqsAWS(
	awsConnection *aws_adapter.AWSConnection,
) *SqsAWS {
	serviceMessageGroupId := os.Getenv("SQS_MESSAGE_GROUP_ID")
	if serviceMessageGroupId == "" {
		serviceMessageGroupId = "group-id:default"
	}
	return &SqsAWS{
		awsConnection: awsConnection,
		groupId:       serviceMessageGroupId,
	}
}

func (s *SqsAWS) SendMessage(queueUrl string, message string) (err error) {
	decuplication := uuid.New().String()
	ctx := context.Background()

	output, err := s.awsConnection.Sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            &message,
		QueueUrl:               &queueUrl,
		MessageGroupId:         &s.groupId,
		MessageDeduplicationId: aws.String(decuplication),
	})

	if err != nil {
		log.Error(context.Background(), err).Msgf("SQS-SendMessage: Error enviando mensagem %+v para Queue %+v > Error: %+v", message, queueUrl, err.Error())
		return errors.NewError(errors.QueueMessageError, err)
	}

	log.Info(ctx).Msgf("SQS-SendMessage: Mensagem com ID: %+v enviada com sucesso", *output.MessageId)

	return nil
}

func (s *SqsAWS) SendMessageWithContext(ctx context.Context, queueUrl string, message string) (err error) {
	decuplication := uuid.New().String()
	output, err := s.awsConnection.Sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            &message,
		QueueUrl:               &queueUrl,
		MessageGroupId:         &s.groupId,
		MessageDeduplicationId: aws.String(decuplication),
	})

	if err != nil {
		log.Error(ctx, err).Msgf("SQS-SendMessageWithContext: Error enviando mensagem %+v para Queue %+v > Error: %+v", message, queueUrl, err.Error())
		return errors.NewError(errors.QueueMessageError, err)
	}

	log.Info(ctx).Msgf("SQS-SendMessageWithContext: Mensagem com ID: %+v enviada com sucesso", *output.MessageId)

	return nil
}
