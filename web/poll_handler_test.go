package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePollBadJSON(t *testing.T) {
	const brokenJSON = `{
    "nope":
  }`
	body := strings.NewReader(brokenJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	handler := NewPollHandler(pollCacheStub, pollStorageStub)
	handler.CreatePoll(res, req)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreatePollWrongLink(t *testing.T) {
	var bodyJSON = `{
    "survey_id": 123,
    "pre_set_values": {
        "п-1": "в-1",
        "п-2": "в-2",
        "п-3": "в-3"
    }
  }`
	body := strings.NewReader(bodyJSON)
	req, err := http.NewRequest(http.MethodPost, "/polls", body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	handler := NewPollHandler(pollCacheStub, pollStorageStub)
	handler.CreatePoll(res, req)
	want := `http://localhost:8080/polls?survey_id=123&%D0%BF-1=%D0%B2-1&%D0%BF-2=%D0%B2-2&%D0%BF-3=%D0%B2-3`
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, want, string(resBody))
}
