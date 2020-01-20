package v2

import (
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/core/prices"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/qiangxue/fasthttp-routing"
)

type Provider struct {
	Store  *cache.Cache
	Prices *prices.PricesData
}

type controller struct {
	store   *cache.Cache
	service *prices.PricesData
}

func API(g *routing.RouteGroup, p *Provider) {
	c := controller{
		store:   p.Store,
		service: p.Prices,
	}

	g.Get("/info", c.info)
	g.Post("/prices", c.courses)
	g.Get("/<crypto>/<fiat>", c.singleCryptoCourse)
	g.Get("/erc20/<token>/<fiat>", c.singleERC20Course)
	g.Get("/docs/swagger.json", respond.SwaggerJSONHandler(v2))
}
