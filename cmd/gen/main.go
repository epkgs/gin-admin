package main

import (
	"context"
	"fmt"
	"gin-admin/internal/app"
	"gin-admin/internal/config"
	"gin-admin/internal/model/po"
	"gin-admin/internal/types"

	"github.com/spf13/cobra"
	"gorm.io/gen"
)

type Querier interface {
	// Get data by ID
	//
	// where("id=@id")
	Get(id string) (*gen.T, error)
}

func executeGen(app types.AppContext) {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/model/bo",                         //gen代码的输出目录
		ModelPkgPath:  "./internal/model",                            //模型代码的输出目录
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface, //启用默认查询和链式接口
		FieldNullable: true,                                          //允许 Null 的字段生成指针类型
	})

	g.UseDB(app.DB())
	//g.GenerateAllTable() //通过数据中的表生成对应的模型

	fmt.Println("正在应用模型...")
	genModels(g)
	fmt.Println("模型应用完成!")

	// 执行生成
	fmt.Println("正在生成 GORM GEN 代码...")
	g.Execute()
	fmt.Println("GORM GEN 代码生成完成!")
}

func genModels(g *gen.Generator) {
	g.ApplyInterface(
		func(Querier) {},
		po.Logger{},
		po.Menu{},
		po.Role{},
		po.User{},
	)
}

func main() {

	cmd := &cobra.Command{
		Use:   "ginadmin-gen",
		Short: "通过 gorm gen 生成 bo 代码",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()

			configFile, _ := cmd.Flags().GetString("config")
			cfg := config.MustLoad(ctx, configFile)
			app := app.New(ctx, cfg)

			executeGen(app)

			return nil
		},
	}
	cmd.Flags().StringP("config", "c", "config.yaml", "Config file")

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
