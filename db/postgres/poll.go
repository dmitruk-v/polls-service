package postgres

import (
	"context"

	"github.com/dmitruk-v/4service/schema"
	"github.com/jackc/pgx/v5"
)

type PollStorage struct {
	db *pgx.Conn
}

func NewPollStorage(db *pgx.Conn) *PollStorage {
	return &PollStorage{
		db: db,
	}
}

func (stg *PollStorage) InsertPoll(ctx context.Context, poll schema.Poll) error {
	q := `
  INSERT INTO polls
    (survey_id, pre_set_values)
  VALUES
    ($1, $2)
  ON CONFLICT (survey_id) DO
    UPDATE SET pre_set_values=$2
  `
	ctag, err := stg.db.Exec(ctx, q, poll.SurveyID, poll.PreSetValues)
	if err != nil {
		return err
	}
	_ = ctag
	return nil
}

func (stg *PollStorage) GetPollByID(ctx context.Context, id int64) (schema.Poll, error) {
	poll := schema.Poll{}
	return poll, nil
}
