package model

var Models []any = []any{}

func AddMigrationModel(model interface{ TableName() string }) {
	Models = append(Models, model)
}
