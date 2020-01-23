package v2

import (
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/qiangxue/fasthttp-routing"
)

type controller struct {
	store *cache.Cache
}

func API(g *routing.RouteGroup, c *cache.Cache) {
	controller := controller{
		store: c,
	}

	g.Get("/info", controller.info)
	g.Post("/prices", controller.courses)
	g.Get("/<crypto>/<fiat>", controller.singleCryptoCourse)
	g.Get("/erc20/<token>/<fiat>", controller.singleERC20Course)
	g.Get("/docs/swagger.json", respond.SwaggerJSONHandler(v2))
}
