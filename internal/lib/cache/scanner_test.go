package cache_test

import (
	"github.com/maetthu/ngxcpd/internal/lib/cache"
	"reflect"
	"sync"
	"testing"
)

func TestScanDir(t *testing.T) {
	var mutex sync.Mutex
	files := make(map[string]*cache.Entry)

	callback := func(entry *cache.Entry) {
		mutex.Lock()
		files[entry.Key] = entry
		mutex.Unlock()
	}

	if err := cache.ScanDir("../../../testdata/cache_files", callback); err != nil {
		t.Error(err)
	}

	if len(files) != len(CacheFiles) {
		t.Fatal("Incorrect number of cache files returned from ScanDir")
	}

	for _, f := range CacheFiles {
		e, ok := files[f.Key]

		if !ok {
			t.Errorf("Key %s not found in returned results", f.Key)
		}

		if !reflect.DeepEqual(e, f) {
			t.Errorf("Loaded cache metadata (%+v) does not match expected value (%+v)", e, f)
		}
	}

}
