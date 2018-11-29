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
	for _, e := range proxycache.TestdataCacheFiles {
		e.Filename = filepath.Join(testdataDir, e.Filename)
		load, err := proxycache.FromFile(e.Filename)

		if err != nil {
			t.Error(err)
			continue
		}

		if !reflect.DeepEqual(e, load) {
			t.Errorf("Loaded metadata (%+v) does not match expected value (%+v)", load, e)
		}

		h, err := load.Hash()

		if err != nil {
			t.Error(err)
		}

		if filepath.Base(load.Filename) != h {
			t.Errorf("Hash calculated from key does not match hash derived from filename")
		}

		if _, err := load.Response(); err != nil {
			t.Error(err)
		}
	}
}
