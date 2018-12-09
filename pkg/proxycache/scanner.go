package proxycache

import (
	"github.com/karrick/godirwalk"
	"sync"
)

// ScanDir walks a caching directory and calls the callback function for each caching file
func ScanDir(dir string, callback func(*Entry), numWorkers int) error {
	indexQueue := make(chan string, numWorkers*128)
	var indexWg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
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
