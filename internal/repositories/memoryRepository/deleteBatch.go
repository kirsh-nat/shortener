package memoryrepository

import "sync"

func (r *MemoryRepository) DeleteBatch(shortURLs []string, userID string) {
	var wg sync.WaitGroup
	results := make(chan string)

	for _, shortRequest := range shortURLs {
		wg.Add(1)
		go func(shortRequest string) {
			defer wg.Done()
			r.mu.Lock()
			defer r.mu.Unlock()

			if urlData, ok := r.store[shortRequest]; ok && urlData.UserID == userID {
				urlData.Deleted = true
				r.store[shortRequest] = urlData
				results <- shortRequest
			}
		}(shortRequest)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		println("Deleted:", result)
	}
}
