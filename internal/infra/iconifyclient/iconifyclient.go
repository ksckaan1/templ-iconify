package iconifyclient

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-resty/resty/v2"
	"github.com/ksckaan1/templ-iconify/internal/core/customerrors"
	"github.com/samber/lo"
)

type IconifyClient struct {
	client *resty.Client
}

func New() *IconifyClient {
	return &IconifyClient{
		client: resty.New(),
	}
}

const iconURL = "https://api.iconify.design/{prefix}/{icon}.svg"

func (c *IconifyClient) GetIconSVG(ctx context.Context, prefix, icon string) (string, error) {
	resp, err := c.client.NewRequest().
		SetPathParams(map[string]string{
			"prefix": prefix,
			"icon":   icon,
		}).
		SetContext(ctx).
		Get(iconURL)
	if err != nil {
		return "", fmt.Errorf("resty.Get: %w", err)
	}
	if resp.StatusCode() != 200 {
		if resp.StatusCode() == 404 {
			return "", customerrors.ErrIconNotFound
		}
		return "", fmt.Errorf("resty.Get: %d", resp.StatusCode())
	}
	return string(resp.Body()), nil
}

const collectionsURL = "https://api.iconify.design/collections"

func (c *IconifyClient) GetCollections(ctx context.Context) ([]string, error) {
	respBody := map[string]any{}
	resp, err := c.client.NewRequest().
		SetResult(&respBody).
		SetContext(ctx).
		Get(collectionsURL)
	if err != nil {
		return nil, fmt.Errorf("resty.Get: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("resty.Get: %d", resp.StatusCode())
	}
	return lo.MapToSlice(respBody, func(key string, _ any) string {
		return key
	}), nil
}

const collectionIconsURL = "https://api.iconify.design/collection"

type collectionIconsResponse struct {
	Uncategorized []string            `json:"uncategorized"`
	Categories    map[string][]string `json:"categories"`
	Hidden        []string            `json:"hidden"`
}

func (c *IconifyClient) GetCollectionIcons(ctx context.Context, prefix string) ([]string, error) {
	respBody := collectionIconsResponse{}
	resp, err := c.client.NewRequest().
		SetResult(&respBody).
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"prefix": prefix,
		}).
		Get(collectionIconsURL)
	if err != nil {
		return nil, fmt.Errorf("resty.Get: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("resty.Get: %d", resp.StatusCode())
	}

	allIcons := slices.Concat(respBody.Uncategorized, respBody.Hidden)
	for _, category := range respBody.Categories {
		allIcons = slices.Concat(allIcons, category)
	}

	allIcons = lo.Uniq(allIcons)

	return allIcons, nil
}
