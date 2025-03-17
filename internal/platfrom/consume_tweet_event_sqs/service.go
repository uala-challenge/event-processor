package consume_tweet_event_sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	sqs2 "github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	ReceiveMessages(ctx context.Context, queueURL string, batchSize, retries int) ([]types.Message, error)
	DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error
	ProcessMessage(ctx context.Context, message string) error
}

type service struct {
	client *sqs2.Sqs
	log    log.Service
}

type Dependencies struct {
	Client *sqs2.Sqs
	Log    log.Service
}

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
	}
}

func (s *service) ReceiveMessages(ctx context.Context, queueURL string, batchSize, retries int) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: int32(batchSize),
		WaitTimeSeconds:     5,
	}

	resp, err := s.client.Cliente.ReceiveMessage(ctx, input)
	if err != nil {
		s.log.Error(ctx, err, "Error recibiendo mensajes de SQS", nil)
		return nil, err
	}

	return resp.Messages, nil
}

func (s *service) DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	}

	_, err := s.client.Cliente.DeleteMessage(ctx, input)
	if err != nil {
		s.log.Error(ctx, err, "Error eliminando mensaje de SQS", nil)
	}
	return err
}

func (s *service) ProcessMessage(ctx context.Context, message string) error {
	s.log.Info(ctx, "Procesando mensaje:", map[string]interface{}{
		"message": message,
	})
	return nil
}
