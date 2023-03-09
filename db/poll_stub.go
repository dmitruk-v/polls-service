package db

import (
	"context"

	"github.com/dmitruk-v/poll-service/schema"
)

type StubPollStorage struct {
	InsertPollFn  func(ctx context.Context, poll schema.Poll) error
	GetPollByIDFn func(ctx context.Context, surveyID int64) (schema.Poll, error)
}

func NewStubPollStorage() *StubPollStorage {
	return &StubPollStorage{}
}

func (stg *StubPollStorage) InsertPoll(ctx context.Context, poll schema.Poll) error {
	return stg.InsertPollFn(ctx, poll)
}

func (stg *StubPollStorage) GetPollByID(ctx context.Context, surveyID int64) (schema.Poll, error) {
	return stg.GetPollByIDFn(ctx, surveyID)
}
