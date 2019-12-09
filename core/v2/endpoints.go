package v2

import (
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	routing "github.com/qiangxue/fasthttp-routing"
)

type Provider struct {
	Store *cache.Cache
}

type controller struct {
	store *cache.Cache
}

func API(g *routing.RouteGroup, p *Provider) {
	c := controller{
		store: p.Store,
	}

	g.Get("/info", c.info)
	g.Get("/docs/swagger.json", respond.SwaggerJSONHandler(v2))
	g.Post("/prices", c.courses)
}
