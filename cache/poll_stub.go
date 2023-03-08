package cache

import (
	"github.com/dmitruk-v/4service/schema"
)

type StubPollCache struct{}

func NewStubPollCache() *StubPollCache {
	return &StubPollCache{}
}

func (c *StubPollCache) HasSurveyID(id string) (bool, error) {
	return true, nil
}

func (c *StubPollCache) GetPoll(key string) (schema.Poll, error) {
	var poll schema.Poll
	return poll, nil
}

func (c *StubPollCache) SetPoll(poll schema.Poll) error {
	return nil
}
