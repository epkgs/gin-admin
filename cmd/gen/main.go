package main

import (
	"context"
	"fmt"
	"gin-admin/internal/config"
	"gin-admin/internal/model"
	"gin-admin/pkg/gormx"

	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Querier interface {
	// Get data by ID
	//
	// where("id=@id")
	Get(id string) (*gen.T, error)
}

func main() {

	cmd := &cobra.Command{
		Use:   "ginadmin-gen",
		Short: "通过 gorm gen 生成 bo 代码",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()

			configFile, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(ctx, configFile)
			if err != nil {
				return err
			}

			db, err := initDB(cfg.DB)
			if err != nil {
				return err
			}

			executeGen(db)

			return nil
		},
	}
	cmd.Flags().StringP("config", "c", "config.yaml", "Config file")

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func initDB(cfg config.DB) (*gorm.DB, error) {
	resolver := make([]gormx.ResolverConfig, len(cfg.Resolver))
	for i, v := range cfg.Resolver {
		resolver[i] = gormx.ResolverConfig{
			DBType:   v.DBType,
			Sources:  v.Sources,
			Replicas: v.Replicas,
			Tables:   v.Tables,
		}
	}

	return gormx.New(gormx.Config{
		Debug:        cfg.Debug,
		PrepareStmt:  cfg.PrepareStmt,
		DBType:       cfg.Type,
		DSN:          cfg.DSN,
		MaxLifetime:  cfg.MaxLifetime,
		MaxIdleTime:  cfg.MaxIdleTime,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		Resolver:     resolver,
	})
}

func executeGen(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/dao",                              // gen代码的输出目录
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface, // 启用默认查询和链式接口
		FieldNullable: true,                                          // 允许 Null 的字段生成指针类型
	})

	g.UseDB(db)
	//g.GenerateAllTable() //通过数据中的表生成对应的模型

	g.ApplyInterface(
		func(Querier) {},
		model.Models...,
	)

	// 执行生成
	fmt.Println("正在生成 GORM GEN 代码...")
	g.Execute()
	fmt.Println("GORM GEN 代码生成完成!")
}
