package consume_tweet_event_sqs

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tweet_mock "github.com/uala-challenge/event-processor/internal/platfrom/tweet_proccesor/mock"
	"github.com/uala-challenge/event-processor/kit"
	sqs2 "github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	sqs_mock "github.com/uala-challenge/simple-toolkit/pkg/client/sqs/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestReceiveMessages_Error(t *testing.T) {
	mockSQS := sqs_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
	expectedErr := errors.New("Error al recibir mensajes de SQS")

	mockSQS.On("ReceiveMessage", mock.Anything, mock.Anything).
		Return(nil, expectedErr)

	mockLog.On("Error", mock.Anything, expectedErr, "Error recibiendo mensajes de SQS", mock.Anything).Return()

	sd := sqs2.Sqs{Cliente: mockSQS}

	service := NewService(Dependencies{
		Client: &sd,
		Log:    mockLog,
	})

	messages, err := service.ReceiveMessages(context.TODO(), queueURL, 5, 3)
	assert.Nil(t, messages)
	assert.Error(t, err)

	mockSQS.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

// ✅ **Test para DeleteMessage - Éxito**
func TestDeleteMessage_Success(t *testing.T) {
	mockSQS := sqs_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
	receiptHandle := "test-receipt-handle"

	mockSQS.On("DeleteMessage", mock.Anything, mock.Anything).
		Return(&sqs.DeleteMessageOutput{}, nil)

	sd := sqs2.Sqs{Cliente: mockSQS}

	service := NewService(Dependencies{
		Client: &sd,
		Log:    mockLog,
	})

	err := service.DeleteMessage(context.TODO(), queueURL, receiptHandle)
	assert.NoError(t, err)

	mockSQS.AssertExpectations(t)
}

// ❌ **Test para DeleteMessage - Error**
func TestDeleteMessage_Error(t *testing.T) {
	mockSQS := sqs_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
	receiptHandle := "test-receipt-handle"
	expectedErr := errors.New("Error al eliminar mensaje de SQS")

	mockSQS.On("DeleteMessage", mock.Anything, mock.Anything).
		Return(nil, expectedErr)

	mockLog.On("Error", mock.Anything, expectedErr, "Error eliminando mensaje de SQS", mock.Anything).Return()

	sd := sqs2.Sqs{Cliente: mockSQS}

	service := NewService(Dependencies{
		Client: &sd,
		Log:    mockLog,
	})

	err := service.DeleteMessage(context.TODO(), queueURL, receiptHandle)
	assert.Error(t, err)

	mockSQS.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

// ✅ **Test para ProcessMessage - Éxito**
func TestProcessMessage_Success(t *testing.T) {
	mockProcessor := tweet_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet-123",
		UserID:  "user-456",
		Created: 1710087745,
	}
	tweetJSON, _ := json.Marshal(tweet)

	mockProcessor.On("Accept", mock.Anything, tweet).Return(nil)

	service := NewService(Dependencies{
		Processor: mockProcessor,
		Log:       mockLog,
	})

	err := service.ProcessMessage(context.TODO(), string(tweetJSON))
	assert.NoError(t, err)

	mockProcessor.AssertExpectations(t)
}

// ❌ **Test para ProcessMessage - Error al deserializar JSON**
func TestProcessMessage_JSONError(t *testing.T) {
	mockProcessor := tweet_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	invalidJSON := `{"tweet_id":"tweet-123","user_id":"user-456","created":INVALID_NUMBER}`

	expectedErr := errors.New("Error deserializando mensaje")
	mockLog.On("WrapError", mock.Anything, "Error deserializando mensaje").Return(expectedErr)

	service := NewService(Dependencies{
		Processor: mockProcessor,
		Log:       mockLog,
	})

	err := service.ProcessMessage(context.TODO(), invalidJSON)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockLog.AssertExpectations(t)
}

// ❌ **Test para ProcessMessage - Error al procesar el tweet**
func TestProcessMessage_ProcessingError(t *testing.T) {
	mockProcessor := tweet_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet-123",
		UserID:  "user-456",
		Created: 1710087745,
	}
	tweetJSON, _ := json.Marshal(tweet)

	expectedErr := errors.New("Error procesando mensaje")
	mockProcessor.On("Accept", mock.Anything, tweet).Return(expectedErr)
	mockLog.On("WrapError", mock.Anything, "Error procesando mensaje").Return(expectedErr)

	service := NewService(Dependencies{
		Processor: mockProcessor,
		Log:       mockLog,
	})

	err := service.ProcessMessage(context.TODO(), string(tweetJSON))
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockProcessor.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}
