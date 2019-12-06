package v2

import (
	"github.com/button-tech/utils-price-tool/pkg/storage"
	routing "github.com/qiangxue/fasthttp-routing"
)

type Provider struct {
	Store *storage.Cache
}

type controller struct {
	store *storage.Cache
}

func API(g *routing.RouteGroup, p *Provider) {
	c := controller{
		store: p.Store,
	}

	g.Post("/prices", c.courses)
	g.Get("/info", c.info)
}
