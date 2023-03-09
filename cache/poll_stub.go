package cache

import (
	"github.com/dmitruk-v/poll-service/schema"
)

type StubPollCache struct {
	GetPollFn func(surveyID int64) (schema.Poll, error)
	SetPollFn func(poll schema.Poll) error
}

func NewStubPollCache() *StubPollCache {
	return &StubPollCache{}
}

func (c *StubPollCache) GetPoll(surveyID int64) (schema.Poll, error) {
	return c.GetPollFn(surveyID)
}

func (c *StubPollCache) SetPoll(poll schema.Poll) error {
	return c.SetPollFn(poll)
}
