package postgres

import (
	"context"
	"errors"

	"github.com/dmitruk-v/poll-service/schema"
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
	_, err := stg.db.Exec(ctx, q, poll.SurveyID, poll.PreSetValues)
	if err != nil {
		return err
	}
	return nil
}

func (stg *PollStorage) GetPollByID(ctx context.Context, surveyID int64) (schema.Poll, error) {
	q := `
  SELECT
    survey_id, pre_set_values
  FROM polls
  WHERE survey_id=$1
  `
	var poll schema.Poll
	row := stg.db.QueryRow(ctx, q, surveyID)
	if err := row.Scan(&poll.SurveyID, &poll.PreSetValues); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return poll, schema.NewErrPollNotFound(surveyID)
		}
		return poll, err
	}
	return poll, nil
}
