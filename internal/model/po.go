package model

type Entity interface {
	TableName() string
}

func Models() []any {
	return []any{
		&User{},
		&Role{},
		&Menu{},
		&Logger{},
	}
}
