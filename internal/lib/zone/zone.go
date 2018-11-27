package zone

import (
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	"github.com/patrickmn/go-cache"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Zone is responsible for maintaining a cache file index for a single zone
type Zone struct {
	Path  string
	Cache *cache.Cache
}

// Warmup scans the cache directory and adds its contents to the indexer queue
func (z *Zone) Warmup(numWorkers int) error {
	return proxycache.ScanDir(z.Path, func(entry *proxycache.Entry) {
		if h, err := entry.Hash(); err == nil {
			z.Cache.Set(h, entry, time.Until(entry.Expire))
		}
	}, numWorkers)
}

// Watch starts listening for filesystem changes in cache directory
func (z *Zone) Watch() error {
	c := make(chan notify.EventInfo, 100)

	watchFor := []notify.Event{
		notify.InMovedTo,
		notify.Remove,
	}

	if err := notify.Watch(filepath.Join(z.Path, "..."), c, watchFor...); err != nil {
		return err
	}

	defer notify.Stop(c)

	for e := range c {
		switch e.Event() {
		case notify.InMovedTo:
			if ce, err := proxycache.FromFile(e.Path()); err == nil {
				if h, err := ce.Hash(); err == nil {
					z.Cache.Set(h, ce, time.Until(ce.Expire))
				}
			}
		case notify.Remove:
			// there may be false positives for temporary files created by nginx before moving it to the
			// final destination. since its file name isn't a valid hash, just ignore it.
			f := filepath.Base(e.Path())
			z.Cache.Delete(f)
		}

		log.Println("Got event:", e)
		log.Printf("%+v\n", e.Path())
		log.Printf("%+v\n", e.Sys())
	}

	return nil
}

// Delete removes an entry from cache and from filesystem
func (z *Zone) Delete(h string) {
	if e, ok := z.Cache.Get(h); ok {
		f := e.(*proxycache.Entry).Filename
		_ = os.Remove(f)
		z.Cache.Delete(h)
	}
}

// WalkNDelete calls function for each entry in cache and removes it if func returns true
func (z *Zone) WalkNDelete(filter func(entry *proxycache.Entry) bool) {
	// TODO: Items() copies *whole* cache into a new map... whhich doesn't sound particularly efficient
	for k, v := range z.Cache.Items() {
		if filter(v.Object.(*proxycache.Entry)) {
			z.Delete(k)
		}
	}
}

// NewZone creates a new indexer instance for given path
func NewZone(path string) (*Zone, error) {
	c := cache.New(cache.NoExpiration, 5*time.Minute)

	return &Zone{Path: path, Cache: c}, nil
}
