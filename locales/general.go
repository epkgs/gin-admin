package locales

import "github.com/epkgs/i18n"

var General = i18n.NewBundle("general")

func init() {
	General.Load()
}
