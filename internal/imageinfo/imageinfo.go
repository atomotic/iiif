package imageinfo

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImageInfo holds information about a local image file
type ImageInfo struct {
	Path   string
	Width  int
	Height int
	Format string
}

// GetImageInfo reads an image file and returns its dimensions and format
func GetImageInfo(path string) (*ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config: %w", err)
	}

	// Determine MIME type from format
	mimeType := "image/" + format
	if format == "jpeg" {
		mimeType = "image/jpeg"
	} else if format == "png" {
		mimeType = "image/png"
	} else if format == "webp" {
		mimeType = "image/webp"
	}

	return &ImageInfo{
		Path:   path,
		Width:  config.Width,
		Height: config.Height,
		Format: mimeType,
	}, nil
}

// ScanDirectory scans a directory for image files and returns their info
func ScanDirectory(dir string) ([]*ImageInfo, error) {
	var images []*ImageInfo

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		info, err := GetImageInfo(path)
		if err != nil {
			// Skip files that can't be decoded
			fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", entry.Name(), err)
			continue
		}

		images = append(images, info)
	}

	return images, nil
}
