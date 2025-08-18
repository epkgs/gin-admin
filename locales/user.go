package locales

import "github.com/epkgs/i18n"

var User = i18n.NewBundle("user")

func init() {
	User.Load()
}
