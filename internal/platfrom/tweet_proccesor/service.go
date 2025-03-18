package tweet_proccesor

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/event-processor/kit"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *redis.Client
	log    log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
	}
}

func (s service) Accept(ctx context.Context, tweet kit.Tweet) error {
	followers, err := s.client.ZRange(ctx, "followers:"+tweet.UserID, 0, -1).Result()
	if err != nil {
		return s.log.WrapError(err, "Error obteniendo seguidores")
	}

	for _, follower := range followers {
		err := s.processes(ctx, tweet, follower)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s service) processes(ctx context.Context, tweet kit.Tweet, follower string) error {
	_, err := s.client.HSet(ctx, tweet.TweetID, map[string]interface{}{
		"tweet_id":   tweet.TweetID,
		"user_id":    tweet.UserID,
		"created_at": tweet.Created,
	}).Result()

	if err != nil {
		return s.log.WrapError(err, "Error al guardar el item")
	}

	_, err = s.client.ZAdd(ctx, fmt.Sprintf("timeline:%s", follower), redis.Z{
		Score:  float64(tweet.Created),
		Member: tweet.TweetID,
	}).Result()

	if err != nil {
		return s.log.WrapError(err, "Error al guardar el item")
	}
	return nil
}
