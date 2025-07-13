package imageapi

import (
	"context"
	"pantheon-auth/graph/model"
)

type API interface {
	SearchSingleImage(ctx context.Context, query string) (*model.Image, error)
}
