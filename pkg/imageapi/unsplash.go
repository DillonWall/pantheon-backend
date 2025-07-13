package imageapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"pantheon-auth/graph/model"
)

type UnsplashClient struct {
	apiKey string
	client *http.Client
}

func NewUnsplashClient() *UnsplashClient {
	return &UnsplashClient{
		apiKey: os.Getenv("UNSPLASH_ACCESS_KEY"),
		client: &http.Client{},
	}
}

type UnsplashSearchResult struct {
	Results []struct {
		ID   string `json:"id"`
		URLs struct {
			Regular string `json:"regular"`
			Thumb   string `json:"thumb"`
		} `json:"urls"`
		Description string `json:"description"`
	} `json:"results"`
}

type UnsplashPhotoResult struct {
	ID   string `json:"id"`
	Tags []struct {
		Title string `json:"title"`
	}
}

func (c *UnsplashClient) SearchSingleImage(ctx context.Context, query string) (*model.Image, error) {
	// GET first photo from search
    encodedQuery := url.QueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.unsplash.com/search/photos?query=%s&client_id=%s&per_page=%s",
			encodedQuery, c.apiKey, "1"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResult UnsplashSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	if len(searchResult.Results) == 0 {
		return nil, fmt.Errorf("no images found for query: %s (API returned 0 results)", query)
	}
	r := searchResult.Results[0]

	// GET photo by ID and attach tags
	req, err = http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.unsplash.com/photos/%s?client_id=%s",
			r.ID, c.apiKey), nil)
	if err != nil {
		return nil, err
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var photoResult UnsplashPhotoResult
	if err := json.NewDecoder(resp.Body).Decode(&photoResult); err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	for _, tag := range photoResult.Tags {
		tags = append(tags, tag.Title)
	}

	// Return object
	image := &model.Image{
		ImageID:    r.ID,
		Thumbnails: r.URLs.Thumb,
		Preview:    r.URLs.Regular,
		Title:      r.Description,
		Source:     model.ImageSourceUnsplash,
		Tags:       tags,
	}

	return image, nil
}
