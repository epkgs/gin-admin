package static

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Option struct {
	ExcludedPathPrefixes []string
	Root                 string
}

type WithOption func(*Option)

func New(opts ...WithOption) gin.HandlerFunc {

	opt := &Option{
		ExcludedPathPrefixes: []string{},
		Root:                 "",
	}

	for _, fn := range opts {
		fn(opt)
	}

	return func(c *gin.Context) {
		if hasPrefix(c.Request.URL.Path, opt.ExcludedPathPrefixes...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(opt.Root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(opt.Root, "index.html")
		}
		c.File(fpath)
		c.Abort()
	}
}

func hasPrefix(path string, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return false
	}

	pathLen := len(path)
	for _, p := range prefixes {

		if pl := len(p); pathLen >= pl && path[:pl] == p {
			return true
		}
	}
	return false
}
