package postmark

import (
	postmarkc "github.com/keighl/postmark"
)

type Config struct {
	AccountKey string
}

func (config *Config) Client() *postmarkc.Client {
	return postmarkc.NewClient("", config.AccountKey)
}
