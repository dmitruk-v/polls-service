package web

import (
	"os"
	"testing"

	"github.com/dmitruk-v/poll-service/cache"
	"github.com/dmitruk-v/poll-service/db"
)

var (
	pollCacheStub    *cache.StubPollCache
	pollStorageStub  *db.StubPollStorage
	htmlRendererStub *StubHTMLRenderer
	pollHandler      *PollHandler
)

func TestMain(m *testing.M) {
	pollCacheStub = cache.NewStubPollCache()
	pollStorageStub = db.NewStubPollStorage()
	htmlRendererStub = NewStubHTMLRender()
	pollHandler = NewPollHandler(pollCacheStub, pollStorageStub, htmlRendererStub)
	os.Exit(m.Run())
}
