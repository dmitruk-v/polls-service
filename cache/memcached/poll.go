package memcached

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dmitruk-v/4service/schema"
)

type PollCache struct {
	client *memcache.Client
}

func NewPollCache(client *memcache.Client) *PollCache {
	return &PollCache{
		client: client,
	}
}

func (c *PollCache) GetPoll(key string) (schema.Poll, error) {
	var poll schema.Poll
	item, err := c.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return poll, fmt.Errorf("key not found: %w", err)
		}
		return poll, err
	}
	surveyID, err := strconv.ParseInt(item.Key, 10, 64)
	if err != nil {
		return poll, err
	}
	poll.SurveyID = surveyID
	if err := json.Unmarshal(item.Value, &poll.PreSetValues); err != nil {
		return poll, err
	}
	return poll, nil
}

func (c *PollCache) SetPoll(poll schema.Poll) error {
	preVals, err := json.Marshal(poll.PreSetValues)
	if err != nil {
		return err
	}
	// TODO: set expiration from ENV variable
	key := fmt.Sprintf("%v", poll.SurveyID)
	item := &memcache.Item{Key: key, Value: preVals, Expiration: 24 * 60 * 60}
	if err := c.client.Set(item); err != nil {
		return err
	}
	return nil
}
