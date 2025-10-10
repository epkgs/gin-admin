package v1

import (
	"gin-admin/internal/dto"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type API struct {
	app    types.AppContext
	engine *gin.Engine
}

var _ types.RoutableHandler = (*API)(nil)

func NewAPI(app types.AppContext) *API {
	return &API{
		app: app,
	}
}

func (a *API) RegRoutes(group *gin.RouterGroup, engine *gin.Engine) {

	g := group.Group("apis")
	g.Use(
		a.app.Middlewares().Auth(),
	)

	a.engine = engine

	g.GET("", a.List)
}

func (a *API) List(c *gin.Context) {
	routes := a.engine.Routes()

	type RouteInfo struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	}

	routeInfos := make([]RouteInfo, len(routes))
	for i, route := range routes {
		routeInfos[i] = RouteInfo{
			Method: route.Method,
			Path:   route.Path,
		}
	}

	response.List(c, routeInfos, &dto.Pager{
		Total: int64(len(routeInfos)),
	})

}
