package locales

import "github.com/epkgs/i18n"

var (
	Def  = i18n.NewBundle("default")
	DB   = i18n.NewBundle("database")
	Http = i18n.NewBundle("http")
	User = i18n.NewBundle("user")
)
