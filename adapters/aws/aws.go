package aws

import (
	"context"

	"github.com/adhfoundation/layer-tools/middlewares"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go/middleware"
)

type AWSConnection struct {
	Sqs *sqs.Client
}

func NewAWSConnection() (*AWSConnection, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		return stack.Serialize.Insert(middleware.SerializeMiddlewareFunc("AwsSqsLogMessageMiddleware", middlewares.AwsSqsLogMessageMiddleware), "OperationSerializer", middleware.After)
	})

	return &AWSConnection{
		Sqs: sqs.NewFromConfig(cfg),
	}, nil
}

func NewAWSConnectionWithContext(ctx context.Context) (*AWSConnection, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &AWSConnection{
		Sqs: sqs.NewFromConfig(cfg),
	}, nil
}
