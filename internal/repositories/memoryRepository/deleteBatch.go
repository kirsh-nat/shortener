package memoryrepository

func (r *MemoryRepository) DeleteBatch(shortURLs []string, userID string) {
	for _, shortRequest := range shortURLs {
		if _, ok := r.store[shortRequest]; !ok {
			continue
		}
		if r.store[shortRequest].UserID != userID {
			continue
		}

		urlData := r.store[shortRequest]
		urlData.Deleted = true
		r.store[shortRequest] = urlData
	}
}
