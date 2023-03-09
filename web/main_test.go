package web

import (
	"os"
	"testing"

	"github.com/dmitruk-v/poll-service/cache"
	"github.com/dmitruk-v/poll-service/db"
)

var (
	pollCacheStub   *cache.StubPollCache
	pollStorageStub *db.StubPollStorage
)

func TestMain(m *testing.M) {
	pollCacheStub = cache.NewStubPollCache()
	pollStorageStub = db.NewStubPollStorage()
	os.Exit(m.Run())
}
