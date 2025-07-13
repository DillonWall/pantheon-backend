package imageapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"pantheon-auth/graph/model"
	"strings"
)

type PixabayClient struct {
	apiKey string
	client *http.Client
}

func NewPixabayClient() *PixabayClient {
	return &PixabayClient{
		apiKey: os.Getenv("PIXABAY_API_KEY"),
		client: &http.Client{},
	}
}

type PixabaySearchResult struct {
	Hits []struct {
		ID           json.Number `json:"id"`
		PreviewURL   string      `json:"previewURL"`
		WebformatURL string      `json:"webformatURL"`
		Tags         string      `json:"tags"`
	} `json:"hits"`
}

func (c *PixabayClient) SearchSingleImage(ctx context.Context, query string) (*model.Image, error) {
	// GET first photo from search
	encodedQuery := url.QueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://pixabay.com/api?q=%s&key=%s&per_page=%s",
			encodedQuery, c.apiKey, "3"), nil) // Note: per_page min is 3
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResult PixabaySearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	if len(searchResult.Hits) == 0 {
		return nil, fmt.Errorf("no images found for query: %s (API returned 0 results)", query)
	}
	r := searchResult.Hits[0]

	// Process tags
	tags := strings.Split(r.Tags, ", ")

	// Return object
	image := &model.Image{
		ImageID:    string(r.ID),
		Thumbnails: r.PreviewURL,   // PreviewURL is the smallest
		Preview:    r.WebformatURL, // Going with WebformatURL since it seems more like a "preview"
		Source:     model.ImageSourcePixabay,
		Tags:       tags,
	}

	return image, nil
}
