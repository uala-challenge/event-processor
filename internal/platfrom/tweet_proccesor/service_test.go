package tweet_proccesor

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uala-challenge/event-processor/kit"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestAccept_ErrorGettingFollowers(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet-123",
		UserID:  "user-456",
		Created: 1710087745,
	}

	expectedErr := errors.New("Redis ZRange error")
	mk.ExpectZRange("followers:user-456", 0, -1).SetErr(expectedErr)

	mockLog.On("WrapError", mock.Anything, "Error obteniendo seguidores").Return(expectedErr)

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	err := service.Accept(context.TODO(), tweet)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	err = mk.ExpectationsWereMet()
	assert.NoError(t, err)
	mockLog.AssertExpectations(t)
}

func TestAccept_ErrorProcessingTweet(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet-123",
		UserID:  "user-456",
		Created: 1710087745,
	}

	expectedErr := errors.New("Error al guardar el item")
	mk.ExpectZRange("followers:user-456", 0, -1).SetVal([]string{"follower-1"})
	mk.ExpectHSet("tweet-123", map[string]interface{}{
		"tweet_id":   "tweet-123",
		"user_id":    "user-456",
		"created_at": int64(1710087745),
	}).SetErr(expectedErr)

	mockLog.On("WrapError", mock.Anything, "Error al guardar el item").Return(expectedErr)

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	err := service.Accept(context.TODO(), tweet)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	err = mk.ExpectationsWereMet()
	assert.NoError(t, err)
	mockLog.AssertExpectations(t)
}

func TestAccept_Success(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet-123",
		UserID:  "user-456",
		Created: 1710087745,
	}

	mk.ExpectZRange("followers:user-456", 0, -1).SetVal([]string{"follower-1"})

	mk.ExpectHSet("tweet-123", map[string]interface{}{
		"tweet_id":   "tweet-123",
		"user_id":    "user-456",
		"created_at": int64(1710087745),
	}).SetVal(int64(1))

	mk.ExpectZAdd("timeline:follower-1", redis.Z{
		Score:  float64(tweet.Created),
		Member: tweet.TweetID,
	}).SetVal(int64(1))

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	err := service.Accept(context.TODO(), tweet)

	assert.NoError(t, err)

	err = mk.ExpectationsWereMet()
	assert.NoError(t, err)

	mockLog.AssertExpectations(t)
}
