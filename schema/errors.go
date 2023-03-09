package schema

import "fmt"

type ErrPollNotFound struct {
	surveyID int64
}

func NewErrPollNotFound(surveyID int64) *ErrPollNotFound {
	return &ErrPollNotFound{
		surveyID: surveyID,
	}
}

func (err *ErrPollNotFound) Error() string {
	return fmt.Sprintf("poll with survey_id %v is not found", err.surveyID)
}
