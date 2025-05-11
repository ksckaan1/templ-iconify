package iconifyclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIconSVG(t *testing.T) {
	client := New()
	icon, err := client.GetIconSVG(context.Background(), "mdi", "home")
	require.NoError(t, err)
	t.Log(string(icon))
}

func TestGetCollections(t *testing.T) {
	client := New()
	collections, err := client.GetCollections(context.Background())
	require.NoError(t, err)
	t.Log(collections)
}

func TestGetCollectionIcons(t *testing.T) {
	client := New()
	icons, err := client.GetCollectionIcons(context.Background(), "mdi")
	require.NoError(t, err)
	t.Log(icons)
}
