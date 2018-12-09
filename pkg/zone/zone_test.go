package zone_test

import (
	"github.com/karrick/godirwalk"
	"github.com/maetthu/ngxcpd/internal/pkg/testfixtures"
	"github.com/maetthu/ngxcpd/pkg/proxycache"
	"github.com/maetthu/ngxcpd/pkg/zone"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var testdataDir = "../../../testdata/cache_files"

func initZone(zoneDir string, t *testing.T) (*zone.Zone, func(*testing.T)) {
	tmpdir, err := ioutil.TempDir("", "zone_test")
	src := filepath.Join(testdataDir, zoneDir)

	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("cp", "-a", src, tmpdir)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	z, err := zone.NewZone(filepath.Join(tmpdir, filepath.Base(src)))

	if err != nil {
		t.Fatal(err)
	}

	if err := z.Warmup(runtime.NumCPU()); err != nil {
		t.Fatal(err)
	}

	cleanup := func(t *testing.T) {
		err := os.RemoveAll(tmpdir)

		if err != nil {
			t.Log(err)
		}
	}

	return z, cleanup
}

func TestZone_Warmup(t *testing.T) {
	for zoneDir, cacheFiles := range testfixtures.TestdataCacheFiles {
		z, cleanup := initZone(zoneDir, t)
		defer cleanup(t)

		items := z.Cache.Items()

		if len(items) != len(cacheFiles) {
			t.Error("Number of item in cache should be the same as in the directory")
		}

	items:
		for k, v := range items {
			e := v.Object.(*proxycache.Entry)

			for _, i := range cacheFiles {
				h, _ := i.Hash()

				if e.Key == i.Key && k == h {
					continue items
				}
			}

			t.Errorf("Loaded entry %s matches no test fixture", e.Key)
		}
	}
}

func runWalkNDelete(zoneDir string, t *testing.T, f func(entry *proxycache.Entry) bool) (int, int) {
	z, cleanup := initZone(zoneDir, t)
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
	callback := func(entry *proxycache.Entry) bool { return true }

	for zoneDir := range testfixtures.TestdataCacheFiles {
		filecount, itemcount := runWalkNDelete(zoneDir, t, callback)

		if filecount > 0 {
			t.Error("Test should have deleted all files")
		}

		if itemcount > 0 {
			t.Error("Cache should contain no items anymore at this point")
		}
	}
}

func TestZone_WalkNDelete_Keep(t *testing.T) {
	callback := func(entry *proxycache.Entry) bool { return false }

	for zoneDir, cacheFiles := range testfixtures.TestdataCacheFiles {
		filecount, itemcount := runWalkNDelete(zoneDir, t, callback)

		if filecount != len(cacheFiles) {
			t.Error("Test should not have deleted any files")
		}

		if itemcount != len(cacheFiles) {
			t.Error("Cache should still contain all items")
		}
	}
}

func TestZone_WatchDelete(t *testing.T) {
	for zoneDir, cacheFiles := range testfixtures.TestdataCacheFiles {
		z, cleanup := initZone(zoneDir, t)
		defer cleanup(t)

		tomb, err := z.Watch(4096)

		if err != nil {
			t.Fatal(err)
		}

		for _, f := range cacheFiles {
			path := filepath.Join(z.Path, f.Filename)

			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}

		done := make(chan struct{}, 1)
		timeout := time.After(30 * time.Second)

		go func() {
			for {
				if z.Cache.ItemCount() == 0 {
					done <- struct{}{}
					return
				}

				time.Sleep(100 * time.Millisecond)
			}
		}()

		go func() {
			select {
			case <-done:
				_ = tomb.Killf("killed")
			case <-timeout:
				_ = tomb.Killf("Waiting for inotify to catch up timed out")
			}
		}()

		err = tomb.Wait()

		if err == nil {
			t.Fatal("Expected error from killing Wait()")
		} else if err.Error() != "killed" {
			t.Fatal(err)
		}
	}
}

func TestZone_WatchAdd(t *testing.T) {
	for zoneDir, cacheFiles := range testfixtures.TestdataCacheFiles {
		z, cleanup := initZone(zoneDir, t)
		defer cleanup(t)

		type File struct {
			Filename string
			Content  []byte
		}

		// buffer cache files in memory and remove from disk
		files := make(map[string]File)

		for _, f := range cacheFiles {
			path := filepath.Join(z.Path, f.Filename)

			if cf, err := os.Open(path); err == nil {
				content, err := ioutil.ReadAll(cf)

				if err != nil {
					t.Fatal(err)
				}

				h, _ := f.Hash()

				files[h] = File{Filename: path, Content: content}
				z.Delete(h)
			} else {
				t.Fatal(err)
			}
		}

		if z.Cache.ItemCount() != 0 {
			t.Log("Cache should be clean but isn't")
		}

		// now repopulate cache again
		tomb, err := z.Watch(16384)

		if err != nil {
			t.Fatal(err)
		}

		for _, f := range files {
			tempname := f.Filename + ".42"

			if err := ioutil.WriteFile(tempname, f.Content, os.FileMode(0500)); err != nil {
				t.Fatal(err)
			}

			if err := os.Rename(tempname, f.Filename); err != nil {
				t.Fatal(err)
			}
		}

		done := make(chan struct{}, 1)
		timeout := time.After(30 * time.Second)

		go func() {
			for {
				if z.Cache.ItemCount() == len(cacheFiles) {
					done <- struct{}{}
					return
				}

				time.Sleep(100 * time.Millisecond)
			}
		}()

		go func() {
			select {
			case <-done:
				_ = tomb.Killf("killed")
			case <-timeout:
				_ = tomb.Killf("Waiting for inotify to catch up timed out")
			}
		}()

		err = tomb.Wait()

		if err == nil {
			t.Fatal("Expected error from killing Wait()")
		} else if err.Error() != "killed" {
			t.Fatal(err)
		}
	}
}
