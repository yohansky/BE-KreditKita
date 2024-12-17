package helpers

import (
	"errors"
	"mime/multipart"
	"net/http"
)

func ValidateAndReadFile(file *multipart.FileHeader, maxSize int64, validTypes []string) ([]byte, error) {
	if file.Size > maxSize {
		return nil, errors.New("file size exceeds the limit")
	}

	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		return nil, err
	}

	for _, validType := range validTypes {
		if validType == http.DetectContentType(buffer) {
			return buffer, nil
		}
	}

	return nil, errors.New("invalid file type")
}
