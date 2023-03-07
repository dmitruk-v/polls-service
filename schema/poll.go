package schema

import "context"

type PollStorage interface {
	InsertPoll(ctx context.Context, poll Poll) error
	GetPollByID(ctx context.Context, id int64) (Poll, error)
}

type PollCache interface {
	GetPoll(key string) (Poll, error)
	SetPoll(poll Poll) error
}

type Poll struct {
	SurveyID     int64             `json:"survey_id"`
	PreSetValues map[string]string `json:"pre_set_values"`
}
