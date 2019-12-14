package v2

import (
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/services"
	routing "github.com/qiangxue/fasthttp-routing"
)

type Provider struct {
	Store     *cache.Cache
	GetPrices *services.GetPrices
}

type controller struct {
	store   *cache.Cache
	service *services.GetPrices
}

func API(g *routing.RouteGroup, p *Provider) {
	c := controller{
		store:   p.Store,
		service: p.GetPrices,
	}

	g.Get("/info", c.info)
	g.Post("/prices", c.courses)
	g.Get("/<crypto>/<fiat>", c.singleCryptoCourse)
	g.Get("/erc20/<token>/<fiat>", c.singleERC20Course)
	g.Get("/docs/swagger.json", respond.SwaggerJSONHandler(v2))
}
