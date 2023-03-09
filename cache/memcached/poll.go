package memcached

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dmitruk-v/poll-service/schema"
)

type PollCache struct {
	client *memcache.Client
}

func NewPollCache(client *memcache.Client) *PollCache {
	return &PollCache{
		client: client,
	}
}

func (c *PollCache) GetPoll(surveyID int64) (schema.Poll, error) {
	var poll schema.Poll
	id := strconv.Itoa(int(surveyID))
	item, err := c.client.Get(id)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return poll, schema.NewErrPollNotFound(surveyID)
		}
		return poll, err
	}
	if err := json.Unmarshal(item.Value, &poll); err != nil {
		return poll, err
	}
	return poll, nil
}

func (c *PollCache) SetPoll(poll schema.Poll) error {
	pollBts, err := json.Marshal(poll)
	if err != nil {
		return err
	}
	// TODO: set expiration from ENV variable
	key := strconv.Itoa(int(poll.SurveyID))
	item := &memcache.Item{Key: key, Value: pollBts, Expiration: 24 * 60 * 60}
	if err := c.client.Set(item); err != nil {
		return err
	}
	return nil
}
