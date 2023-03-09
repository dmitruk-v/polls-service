package schema

import "fmt"

type ErrPollNotFound struct {
	pollID int64
}

func NewErrPollNotFound(pollID int64) *ErrPollNotFound {
	return &ErrPollNotFound{
		pollID: pollID,
	}
}

func (err *ErrPollNotFound) Error() string {
	return fmt.Sprintf("poll with id %v is not found", err.pollID)
}
