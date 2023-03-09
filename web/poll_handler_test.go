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
	"github.com/go-chi/chi/v5"
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
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return nil
	}
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return nil
	}
	body := strings.NewReader(badJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	pollHandler.CreatePoll(res, req)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreatePoll_DBErrorStatus500(t *testing.T) {
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return nil
	}
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return errors.New("db error")
	}
	body := strings.NewReader(goodJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	pollHandler.CreatePoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestCreatePoll_CacheErrorStatus500(t *testing.T) {
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return errors.New("cache error")
	}
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return nil
	}
	body := strings.NewReader(goodJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	pollHandler.CreatePoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestCreatePoll_GoodLink(t *testing.T) {
	pollCacheStub.SetPollFn = func(poll schema.Poll) error {
		return nil
	}
	pollStorageStub.InsertPollFn = func(ctx context.Context, poll schema.Poll) error {
		return nil
	}
	want := `http://localhost:8080/polls/123?%D0%BF-1=%D0%B2-1&%D0%BF-2=%D0%B2-2&%D0%BF-3=%D0%B2-3`
	body := strings.NewReader(goodJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	pollHandler.CreatePoll(res, req)
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, want, string(resBody))
}

func TestGetPoll_BadSurveyID(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/polls/1q?a=1&b=2&c=3", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	// Need to manualy add url params to chi.URLParams
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("survey_id", "1q")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// ----------
	pollHandler.GetPoll(res, req)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestGetPoll_CacheError(t *testing.T) {
	pollCacheStub.GetPollFn = func(surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, errors.New("cache error")
	}
	req, err := http.NewRequest(http.MethodPost, "/polls/123?a=1&b=2&c=3", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	// Need to manualy add url params to chi.URLParams
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("survey_id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// ----------
	pollHandler.GetPoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetPoll_CacheMiss(t *testing.T) {
	pollCacheStub.GetPollFn = func(surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, schema.NewErrPollNotFound(surveyID)
	}
	var fromDB bool
	pollStorageStub.GetPollByIDFn = func(ctx context.Context, surveyID int64) (schema.Poll, error) {
		fromDB = true
		return schema.Poll{}, nil
	}
	req, err := http.NewRequest(http.MethodPost, "/polls/123?a=1&b=2&c=3", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	// Need to manualy add url params to chi.URLParams
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("survey_id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// ----------
	pollHandler.GetPoll(res, req)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, true, fromDB)
}

func TestGetPoll_DBError(t *testing.T) {
	pollCacheStub.GetPollFn = func(surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, nil
	}
	pollStorageStub.GetPollByIDFn = func(ctx context.Context, surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, errors.New("db error")
	}
	req, err := http.NewRequest(http.MethodPost, "/polls/123?a=1&b=2&c=3", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	// Need to manualy add url params to chi.URLParams
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("survey_id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// ----------
	pollHandler.GetPoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetPoll_DBMiss(t *testing.T) {
	pollCacheStub.GetPollFn = func(surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, schema.NewErrPollNotFound(surveyID)
	}
	pollStorageStub.GetPollByIDFn = func(ctx context.Context, surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, schema.NewErrPollNotFound(surveyID)
	}
	req, err := http.NewRequest(http.MethodPost, "/polls/123?a=1&b=2&c=3", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	// Need to manualy add url params to chi.URLParams
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("survey_id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// ----------
	pollHandler.GetPoll(res, req)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetPoll_OK(t *testing.T) {
	pollCacheStub.GetPollFn = func(surveyID int64) (schema.Poll, error) {
		poll := schema.Poll{
			SurveyID: 123,
			PreSetValues: map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		}
		return poll, nil
	}
	pollStorageStub.GetPollByIDFn = func(ctx context.Context, surveyID int64) (schema.Poll, error) {
		return schema.Poll{}, nil
	}
	req, err := http.NewRequest(http.MethodPost, "/polls/123?a=1&b=2&c=3", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	// Need to manualy add url params to chi.URLParams
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("survey_id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// ----------
	pollHandler.GetPoll(res, req)
	require.Equal(t, http.StatusOK, res.Code)
}
