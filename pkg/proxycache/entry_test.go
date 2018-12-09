package proxycache_test

import (
	"github.com/maetthu/ngxcpd/internal/pkg/testfixtures"
	"github.com/maetthu/ngxcpd/pkg/proxycache"
	"path/filepath"
	"reflect"
	"testing"
)

var testdataDir = "../../testdata/cache_files"

// TestFromFile checks if the file parser works as expected
func TestFromFile(t *testing.T) {
	for zoneDir, cacheFiles := range testfixtures.TestdataCacheFiles {
		dir := filepath.Join(testdataDir, zoneDir)

		for _, e := range cacheFiles {
			e.Filename = filepath.Join(dir, e.Filename)
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

}
