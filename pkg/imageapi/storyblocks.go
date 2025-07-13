package imageapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"pantheon-auth/graph/model"
	"time"
)

type StoryblocksClient struct {
	apiKey string
	client *http.Client
}

func NewStoryblocksClient() *StoryblocksClient {
	return &StoryblocksClient{
		apiKey: os.Getenv("STORYBLOCKS_API_KEY"),
		client: &http.Client{},
	}
}

type StoryblocksSearchResult struct {
	Results []struct {
		ID           json.Number `json:"id"`
		Title        string      `json:"title"`
		ThumbnailURL string      `json:"thumbnail_url"`
		PreviewURL   string      `json:"preview_url"`
		Keywords     []string    `json:"keywords"`
	} `json:"results"`
}

func (c *StoryblocksClient) SearchSingleImage(ctx context.Context, query string) (*model.Image, error) {
	// GET first photo from search
	expires := time.Now().Unix() + 100
	resource := "/api/v2/images/search"
	secretKey := os.Getenv("STORYBLOCKS_SECRET_KEY")
	projectId := os.Getenv("STORYBLOCKS_PROJECT_ID")
	userId := os.Getenv("STORYBLOCKS_USER_ID")

	hmacKey := []byte(secretKey + fmt.Sprintf("%d", expires))
	h := hmac.New(sha256.New, hmacKey)
	h.Write([]byte(resource))
	hmacDigest := hex.EncodeToString(h.Sum(nil))

	encodedQuery := url.QueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.storyblocks.com/api/v2/images/search?APIKEY=%s&EXPIRES=%d&HMAC=%s&project_id=%s&user_id=%s&keywords=%s&results_per_page=%s&extended=keywords",
			c.apiKey, expires, hmacDigest, projectId, userId, encodedQuery, "1"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResult StoryblocksSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	if len(searchResult.Results) == 0 {
		return nil, fmt.Errorf("no images found for query: %s (API returned 0 results)", query)
	}
	r := searchResult.Results[0]

	// Return object
	image := &model.Image{
		ImageID:    string(r.ID),
		Title:      r.Title,
		Preview:    r.PreviewURL,
		Thumbnails: r.ThumbnailURL,
		Source:     model.ImageSourceStoryblocks,
		Tags:       r.Keywords,
	}

	return image, nil
}

