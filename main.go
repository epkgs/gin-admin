package main

import (
	"gin-admin/cmd"

	"github.com/spf13/cobra"
)

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

	rootCmd.AddCommand(cmd.StartCmd())
	rootCmd.AddCommand(cmd.StopCmd())
	rootCmd.AddCommand(cmd.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
