package config

import (
	"github.com/uala-challenge/tweets-service/internal/retrieve_tweet"
	"github.com/uala-challenge/tweets-service/internal/store_tweet"
)

type UsesCasesConfig struct {
	Store    store_tweet.Config    `json:"store"`
	Retrieve retrieve_tweet.Config `json:"retrieve"`
}
