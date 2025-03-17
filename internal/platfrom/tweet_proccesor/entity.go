package tweet_proccesor

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/event-processor/kit"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply(ctx context.Context, tweet kit.Tweet) (*resty.Response, error)
}

type Dependencies struct {
	Client *redis.Client
	Log    log.Service
}
