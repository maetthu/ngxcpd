package proxycache_test

import (
	"github.com/maetthu/ngxcpd/internal/pkg/testfixtures"
	"github.com/maetthu/ngxcpd/pkg/proxycache"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestScanDir(t *testing.T) {
	for zones, cacheFiles := range testfixtures.TestdataCacheFiles {
		dir := filepath.Join(testdataDir, zones)

		var mutex sync.Mutex
		files := make(map[string]*proxycache.Entry)

		callback := func(entry *proxycache.Entry) {
			mutex.Lock()
			files[entry.Key] = entry
			mutex.Unlock()
		}

		if err := proxycache.ScanDir(dir, callback, runtime.NumCPU()); err != nil {
			t.Error(err)
		}

		if len(files) != len(cacheFiles) {
			t.Fatal("Incorrect number of cache files returned from ScanDir")
		}

		for _, f := range cacheFiles {
			e, ok := files[f.Key]

			if !ok {
				t.Errorf("Key %s not found in returned results", f.Key)
			}

			if !reflect.DeepEqual(e, f) {
				t.Errorf("Loaded cache metadata (%+v) does not match expected value (%+v)", e, f)
			}
		}
	}

}
