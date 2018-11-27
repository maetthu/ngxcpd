package proxycache_test

import (
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	"path/filepath"
	"reflect"
	"testing"
)

var testdataDir = "../../../testdata/cache_files"

// TestFromFile checks if the file parser works as expected
func TestFromFile(t *testing.T) {
	for _, e := range proxycache.Test_CacheFiles {
		e.Filename = filepath.Join(testdataDir, e.Filename)
		load, err := proxycache.FromFile(e.Filename)

		if err != nil {
			t.Error(err)
			continue
		}

		if !reflect.DeepEqual(e, load) {
			t.Errorf("Loaded cache metadata (%+v) does not match expected value (%+v)", load, e)
		}

		h, err := e.Hash()

		if err != nil {
			t.Error(err)
		}

		if filepath.Base(e.Filename) != h {
			t.Errorf("Hash calculated from key does not match hash derived from filename")
		}
	}
}
