package proxycache_test

import (
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// CacheFiles contains expected metadata of files in testdata/cache_files
var CacheFiles = []*proxycache.Entry{
	{
		Filename:     "testdata/cache_files/3/f9/f46162e692a012e2b97c2cc1d2c33f93",
		Version:      5,
		Expire:       time.Unix(1541153027, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151227, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?4",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/d/c0/6015a2835a7797df66ca5dd689541c0d",
		Version:      5,
		Expire:       time.Unix(1541153026, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151226, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?2",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/d/54/8b38764aa969ab8edac1ecac0dfee54d",
		Version:      5,
		Expire:       time.Unix(1541153029, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151229, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?9",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/6/54/bec21cc16bcb756a28bf0b1b72d90546",
		Version:      5,
		Expire:       time.Unix(1541153027, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151227, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?3",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/5/93/9902967e49a31ec37c8d00d64ffaf935",
		Version:      5,
		Expire:       time.Unix(1541153028, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151228, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?7",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/9/d2/9c0b399e0c510e4eb087a0a6369e7d29",
		Version:      5,
		Expire:       time.Unix(1541153026, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151226, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?1",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/1/f9/577a88d09656e8c0b430a472c33fcf91",
		Version:      5,
		Expire:       time.Unix(1541153029, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151229, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?10",
		HeaderStart:  373,
		BodyStart:    603,
	},
	{
		Filename:     "testdata/cache_files/1/00/273725162d80a5f4a63c9c70caf5e001",
		Version:      5,
		Expire:       time.Unix(1541153028, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151228, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?6",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/1/29/b6dbbb30efda04920795036baba9e291",
		Version:      5,
		Expire:       time.Unix(1541153029, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151229, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?8",
		HeaderStart:  372,
		BodyStart:    602,
	},
	{
		Filename:     "testdata/cache_files/b/4f/c71483bd3e6186f4e1c4d4ca4b2ba4fb",
		Version:      5,
		Expire:       time.Unix(1541153027, 0),
		LastModified: time.Unix(1541149732, 0),
		Date:         time.Unix(1541151227, 0),
		Etag:         "\"5bdc1424-6\"",
		Key:          "http://localhost/ngxcpd.txt?5",
		HeaderStart:  372,
		BodyStart:    602,
	},
}

// TestFromFile checks if the file parser works as expected
func TestFromFile(t *testing.T) {
	for _, e := range CacheFiles {
		e.Filename = filepath.Join("../../..", e.Filename)
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
