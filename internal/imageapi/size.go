package imageapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Image struct {
	Context  string `json:"@context"`
	Protocol string `json:"protocol"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Sizes    []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"sizes"`
	Tiles []struct {
		Width        int   `json:"width"`
		Height       int   `json:"height"`
		ScaleFactors []int `json:"scaleFactors"`
	} `json:"tiles"`
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Profile        string   `json:"profile"`
	MaxWidth       int      `json:"maxWidth"`
	MaxHeight      int      `json:"maxHeight"`
	ExtraQualities []string `json:"extraQualities"`
	ExtraFeatures  []string `json:"extraFeatures"`
}

func GetSize(imageapi string) ([]int, error) {
	var image Image

	resp, err := http.Get(imageapi)
	if err != nil {
		return nil, fmt.Errorf("error fetching image info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting image size: HTTP %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		return nil, fmt.Errorf("error decoding image info: %v", err)
	}

	return []int{image.Width, image.Height}, nil
}
