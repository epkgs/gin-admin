package routers

import (
	"context"
	"gin-admin/internal/types"

	"github.com/gin-gonic/gin"
)

type Routers struct {
	app         types.AppContext
	routes      map[string]types.RouteRegister
	middlewares types.Middlewares
}

var _ types.Routers = (*Routers)(nil)

func NewRouters(app types.AppContext) *Routers {
	return &Routers{
		app:         app,
		routes:      map[string]types.RouteRegister{},
		middlewares: NewMiddlewares(app),
	}
}

func (r *Routers) ApiGroup(prefix string, register types.RouteRegister) {
	r.middlewares.Auth().Include(prefix)
	r.middlewares.Casbin().Include(prefix)
	r.middlewares.CopyBody().Include(prefix)
	r.middlewares.Logger().Include(prefix)
	r.middlewares.Trace().Include(prefix)
	r.middlewares.RateLimiter().Include(prefix)

	r.routes[prefix] = register
}

func (r *Routers) Middlewares() types.Middlewares {
	return r.middlewares
}

func (r *Routers) Init(ctx context.Context, e *gin.Engine) error {

	for prefix, register := range r.routes {
		g := e.Group(prefix)
		if err := register(ctx, g, e); err != nil {
			return err
		}
	}
	return nil
}
