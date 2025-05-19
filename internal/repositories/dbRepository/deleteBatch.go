package dbrepository

import (
	"context"
	"fmt"
	"sync"

	"github.com/kirsh-nat/shortener.git/internal/app"
)

func (r *DBRepository) DeleteBatch(shortURLs []string, userID string) {
	const batchSize = 5
	var wg sync.WaitGroup
	errChan := make(chan error, len(shortURLs))
	urlChan := make(chan string)

	go func() {
		for err := range errChan {
			if err != nil {
				app.Sugar.Info("Error:", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for url := range urlChan {
			if err := r.updateDeletedStatus(url, userID); err != nil {
				errChan <- err
			}
		}
	}()

	for i := 0; i < len(shortURLs); i += batchSize {
		end := i + batchSize
		if end > len(shortURLs) {
			end = len(shortURLs)
		}

		batch := shortURLs[i:end]

		for _, url := range batch {
			urlChan <- url
		}

	}

	close(urlChan)
	wg.Wait()
	close(errChan)
}

func (r *DBRepository) updateDeletedStatus(shortURL, userID string) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx,
		"UPDATE links SET deleted = $1 WHERE short_url = $2 AND user_id = $3", true, shortURL, userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for URL: %s", shortURL)
	}

	return nil
}
