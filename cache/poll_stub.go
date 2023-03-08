package cache

import (
	"github.com/dmitruk-v/4service/schema"
)

type StubPollCache struct {
	HasSurveyFn func(id string) (bool, error)
	GetPollFn   func(key string) (schema.Poll, error)
	SetPollFn   func(poll schema.Poll) error
}

func NewStubPollCache() *StubPollCache {
	return &StubPollCache{}
}

func (c *StubPollCache) HasSurveyID(id string) (bool, error) {
	return c.HasSurveyFn(id)
}

func (c *StubPollCache) GetPoll(id string) (schema.Poll, error) {
	return c.GetPollFn(id)
}

func (c *StubPollCache) SetPoll(poll schema.Poll) error {
	return c.SetPollFn(poll)
}
