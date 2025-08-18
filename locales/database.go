package locales

import "github.com/epkgs/i18n"

var DB = i18n.NewBundle("database")

func init() {
	DB.Load()
}
