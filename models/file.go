package models

import (
	"path/filepath"
	"time"
)

type FileField struct {
	FileName string `orm:"size(300); null"` // Stores the file name
	FilePath string `orm:"size(300); null"` // Stores the file path
}

// Generate a new file name (e.g., unique timestamped name)
func (f *FileField) GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	newName := time.Now().Format("20060102150405") + ext // e.g., 20250121123456.png
	return newName
}

// Save metadata for the uploaded file
func (f *FileField) SaveFileMetadata(fileName, filePath string) {
	f.FileName = fileName
	f.FilePath = filePath
}

// Full URL to access the file (optional, for serving files)
func (f *FileField) GetFileURL(baseURL string) string {
	return baseURL + "/" + f.FilePath
}
