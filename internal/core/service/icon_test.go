package service

import (
	"context"
	"testing"

	"github.com/ksckaan1/templ-iconify/internal/infra/iconifyclient"
	"github.com/stretchr/testify/require"
)

func TestFindIcons(t *testing.T) {
	client := iconifyclient.New()
	app := NewIconService(client)
	ctx := context.Background()
	icons, err := app.FindIcons(ctx, "mdi:home", "mdi:home-outline")
	require.NoError(t, err)
	require.Equal(t, 2, len(icons))
}
