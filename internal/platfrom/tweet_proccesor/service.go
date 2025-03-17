package tweet_proccesor

import (
	"context"
	"github.com/go-resty/resty/v2"
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

	return nil
}
