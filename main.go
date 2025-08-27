package main

import (
	"gin-admin/cmd/app"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/epkgs/i18n/i18ntool extract

// Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "v1.0.0"

// @title Gin Admin
// @description 基于 Gin 的快速启动项目
// @version v1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
func main() {
	rootCmd := &cobra.Command{
		Use:     "ginadmin",
		Short:   "基于 Gin 的快速启动项目",
		Version: VERSION,
	}

	rootCmd.AddCommand(app.StartCmd())
	rootCmd.AddCommand(app.StopCmd())
	rootCmd.AddCommand(app.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
