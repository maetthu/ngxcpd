package cache

import (
	"github.com/karrick/godirwalk"
	"sync"
)

const (
	indexWorkers = 32
)

// ScanDir walks a caching directory and calls the callback function for each caching file
func ScanDir(dir string, callback func(*Entry)) error {
	indexQueue := make(chan string, indexWorkers*128)
	var indexWg sync.WaitGroup

	for i := 0; i < indexWorkers; i++ {
		indexWg.Add(1)

		go func(index chan string) {
			defer indexWg.Done()

			for f := range index {
				if c, err := FromFile(f); err == nil {
					callback(c)
				}
			}
		}(indexQueue)
	}

	err := godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				indexQueue <- osPathname
			}

			return nil
		},
		Unsorted: true,
	})

	close(indexQueue)
	indexWg.Wait()

	return err
}
