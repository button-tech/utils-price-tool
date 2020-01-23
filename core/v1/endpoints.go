package v1

import (
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/qiangxue/fasthttp-routing"
)

type controller struct {
	storage           *cache.Cache
	privateCurrencies map[string][]string
}

func API(g *routing.RouteGroup, c *cache.Cache) {
	controller := controller{
		storage:           c,
		privateCurrencies: privateCurrencies(),
	}

	g.Get("/info", controller.info)
	g.Get("/docs/swagger.json", respond.SwaggerJSONHandler(v1))
	g.Post("/prices", controller.courses)
	g.Post("/change", controller.privatePrices)

}

func privateCurrencies() map[string][]string {
	return map[string][]string{
		"BTC":   {"0x0000000000000000000000000000000000000000", "Bitcoin"},
		"ETH":   {"0x000000000000000000000000000000000000003c", "Ethereum"},
		"ETC":   {"0x000000000000000000000000000000000000003d", "Ethereum Classic"},
		"BCH":   {"0x0000000000000000000000000000000000000091", "Bitcoin Cash"},
		"LTC":   {"0x0000000000000000000000000000000000000002", "Litecoin"},
		"XLM":   {"0x0000000000000000000000000000000000000094", "Stellar"},
		"WAVES": {"0x0000000000000000000000000000000000579bfc", "Waves"},
	}
}
