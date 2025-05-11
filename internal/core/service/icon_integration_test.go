//go:build integration

package service

import (
	"context"
	"testing"

	"github.com/ksckaan1/templ-iconify/internal/core/customerrors"
	"github.com/ksckaan1/templ-iconify/internal/core/domain"
	"github.com/ksckaan1/templ-iconify/internal/infra/iconifyclient"
	"github.com/stretchr/testify/require"
)

func TestIconService_FindIcons(t *testing.T) {
	client := iconifyclient.New()
	app := NewIconService(client)
	ctx := context.Background()

	tests := []struct {
		name string
		args []string
		want func([]*domain.Icon)
		err  require.ErrorAssertionFunc
	}{
		{
			name: "valid with one icon",
			args: []string{"mdi:home"},
			want: func(got []*domain.Icon) {
				require.Equal(t, []*domain.Icon{
					{
						Prefix: "mdi",
						Name:   "home",
					},
				}, got)
			},
			err: require.NoError,
		},
		{
			name: "valid with multiple icons",
			args: []string{"mdi:home", "mdi:home-outline"},
			want: func(got []*domain.Icon) {
				require.Equal(t, []*domain.Icon{
					{
						Prefix: "mdi",
						Name:   "home",
					},
					{
						Prefix: "mdi",
						Name:   "home-outline",
					},
				}, got)
			},
			err: require.NoError,
		},
		{
			name: "icon not found",
			args: []string{"mdi:homesss"},
			want: func(got []*domain.Icon) {
				require.Nil(t, got)
			},
			err: func(tt require.TestingT, err error, i ...any) {
				require.ErrorIs(tt, err, customerrors.ErrIconNotFound)
			},
		},
		{
			name: "icon not found inside multiple icons",
			args: []string{"mdi:homesss", "mdi:home-outline"},
			want: func(got []*domain.Icon) {
				require.Equal(t, []*domain.Icon{
					{
						Prefix: "mdi",
						Name:   "home-outline",
					},
				}, got)
			},
			err: require.NoError,
		},
		{
			name: "invalid icon name",
			args: []string{"invalid"},
			want: func(got []*domain.Icon) {
				require.Nil(t, got)
			},
			err: func(tt require.TestingT, err error, i ...any) {
				require.ErrorIs(tt, err, customerrors.ErrInvalidIconName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icons, err := app.FindIcons(ctx, tt.args...)
			tt.err(t, err)
			tt.want(icons)
		})
	}
}
