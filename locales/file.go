package locales

import "github.com/epkgs/i18n"

var File = i18n.NewBundle("file")

func init() {
	File.Load()
}
