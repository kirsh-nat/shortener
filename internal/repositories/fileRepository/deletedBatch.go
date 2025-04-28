package filerepository

import (
	"encoding/json"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *FileRepository) DeleteBatch(shortURLs []string, userID string) {

	file, err := os.Open(r.filePath)
	if err != nil {
		return
	}
	defer file.Close()

	data := make(map[string]services.UserURLData)
	if err := r.loadData(file, &data); err != nil {
		return
	}

	newData := make(map[string]services.UserURLData)
	for user, userData := range data {
		if userID == user {
			for _, reqShort := range shortURLs {
				if reqShort == userData.Short {
					userData.Deleted = true
				}
			}
			newData[user] = userData
		}
	}

	writeData, err := json.MarshalIndent(newData, "", "   ")
	if err != nil {
		return
	}
	writeData = append(writeData, '\n')

	_, err = file.Write(writeData)
	if err != nil {
		return
	}
}
