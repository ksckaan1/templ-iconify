package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ksckaan1/templ-iconify/internal/core/app"
	"github.com/ksckaan1/templ-iconify/internal/core/service"
	"github.com/ksckaan1/templ-iconify/internal/infra/iconifyclient"
)

func main() {
	ctx := context.Background()
	client := iconifyclient.New()
	iconService := service.NewIconService(client)
	rootCmd := app.NewRootCmd(iconService)

	err := rootCmd.Run(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
