package web

import (
	"os"
	"testing"

	"github.com/dmitruk-v/4service/cache"
	"github.com/dmitruk-v/4service/db"
	"github.com/dmitruk-v/4service/schema"
)

var (
	pollCacheStub   schema.PollCache
	pollStorageStub schema.PollStorage
)

func TestMain(m *testing.M) {
	pollCacheStub = cache.NewStubPollCache()
	pollStorageStub = db.NewStubPollStorage()
	os.Exit(m.Run())
}
