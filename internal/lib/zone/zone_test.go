package zone_test

import (
	"github.com/karrick/godirwalk"
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	"github.com/maetthu/ngxcpd/internal/lib/zone"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var testdataDir = "../../../testdata/cache_files"

func initZone(t *testing.T) (*zone.Zone, func(*testing.T)) {
	dir, err := ioutil.TempDir("", "zone_test")

	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("cp", "-a", testdataDir, dir)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	z, err := zone.NewZone(filepath.Join(dir, filepath.Base(testdataDir)))

	if err != nil {
		t.Fatal(err)
	}

	if err := z.Warmup(runtime.NumCPU()); err != nil {
		t.Fatal(err)
	}

	cleanup := func(t *testing.T) {
		err := os.RemoveAll(dir)

		if err != nil {
			t.Log(err)
		}
	}

	return z, cleanup
}

func TestZone_Warmup(t *testing.T) {
	z, cleanup := initZone(t)
	defer cleanup(t)

	items := z.Cache.Items()

	if len(items) != len(proxycache.TestdataCacheFiles) {
		t.Error("Number of item in cache should be the same as in the directory")
	}

items:
	for k, v := range items {
		e := v.Object.(*proxycache.Entry)

		for _, i := range proxycache.TestdataCacheFiles {
			h, _ := i.Hash()

			if e.Key == i.Key && k == h {
				continue items
			}
		}

		t.Errorf("Loaded entry %s matches no test fixture", e.Key)
	}
}

func runWalkNDelete(t *testing.T, f func(entry *proxycache.Entry) bool) (int, int) {
	z, cleanup := initZone(t)
	defer cleanup(t)

	z.WalkNDelete(f)

	filecount := 0

	err := godirwalk.Walk(z.Path, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				filecount++
			}

			return nil
		},
		Unsorted: true,
	})

	if err != nil {
		t.Error(err)
	}

	return filecount, z.Cache.ItemCount()
}

func TestZone_WalkNDelete_Delete(t *testing.T) {
	f := func(entry *proxycache.Entry) bool { return true }
	filecount, itemcount := runWalkNDelete(t, f)

	if filecount > 0 {
		t.Error("Test should have deleted all files")
	}

	if itemcount > 0 {
		t.Error("Cache should contain no items anymore at this point")
	}
}

func TestZone_WalkNDelete_Keep(t *testing.T) {
	f := func(entry *proxycache.Entry) bool { return false }
	filecount, itemcount := runWalkNDelete(t, f)

	if filecount != len(proxycache.TestdataCacheFiles) {
		t.Error("Test should not have deleted any files")
	}

	if itemcount != len(proxycache.TestdataCacheFiles) {
		t.Error("Cache should still contain all items")
	}
}
