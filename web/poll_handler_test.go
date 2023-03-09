package web

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmitruk-v/poll-service/schema"
	"github.com/stretchr/testify/require"
)

const (
	goodJSON = `{
		"survey_id": 123,
		"pre_set_values": {
				"п-1": "в-1",
				"п-2": "в-2",
				"п-3": "в-3"
		}
	}`
	badJSON = `{
    "nope": "just bad",
		"uh"
  }`
)

func TestCreatePoll_BadJSONStatus400(t *testing.T) {
	body := strings.NewReader(badJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	handler := NewPollHandler(pollCacheStub, pollStorageStub)
	handler.CreatePoll(res, req)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreatePoll_DBErrorStatus500(t *testing.T) {
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return errors.New("db error")
	}
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return nil
	}
	body := strings.NewReader(goodJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	handler := NewPollHandler(pollCacheStub, pollStorageStub)
	handler.CreatePoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestCreatePoll_CacheErrorStatus500(t *testing.T) {
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return nil
	}
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return errors.New("cache error")
	}
	body := strings.NewReader(goodJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	handler := NewPollHandler(pollCacheStub, pollStorageStub)
	handler.CreatePoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestCreatePoll_WrongLink(t *testing.T) {
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return nil
	}
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return nil
	}
	want := `http://localhost:8080/polls?survey_id=0123&%D0%BF-1=%D0%B2-1&%D0%BF-2=%D0%B2-2&%D0%BF-3=%D0%B2-3`
	body := strings.NewReader(goodJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	handler := NewPollHandler(pollCacheStub, pollStorageStub)
	handler.CreatePoll(res, req)
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.Code)
	require.NotEqual(t, want, string(resBody))
}
