package db

import (
	"context"

	"github.com/dmitruk-v/4service/schema"
)

type StubPollStorage struct{}

func NewStubPollStorage() *StubPollStorage {
	return &StubPollStorage{}
}

func (stg *StubPollStorage) InsertPoll(ctx context.Context, poll schema.Poll) error {
	return nil
}

func (stg *StubPollStorage) GetPollByID(ctx context.Context, id int64) (schema.Poll, error) {
	poll := schema.Poll{}
	return poll, nil
}
