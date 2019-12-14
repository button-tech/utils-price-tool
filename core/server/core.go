package server

import (
	"encoding/json"
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	v1 "github.com/button-tech/utils-price-tool/core/v1"
	v2 "github.com/button-tech/utils-price-tool/core/v2"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/services"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"time"
)

type Core struct {
	S   *fasthttp.Server
	R   *routing.Router
	G   *routing.RouteGroup
	Gv2 *routing.RouteGroup

	service *services.GetPrices
	store   *cache.Cache
}

func New(store *cache.Cache, service *services.GetPrices) (c *Core) {
	c = &Core{
		R:       routing.New(),
		store:   store,
		service: service,
	}
	c.R.Use(cors)
	c.initBaseRoute()
	c.fs()

	v1.API(c.G, &v1.Provider{Store: c.store, GetPrices: c.service})
	v2.API(c.Gv2, &v2.Provider{Store: c.store, GetPrices: c.service})
	return
}

func (c *Core) initBaseRoute() {
	c.G = c.R.Group("/courses/v1")
	c.Gv2 = c.R.Group("/courses/v2")
}

func (c *Core) fs() {
	c.S = &fasthttp.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      c.R.HandleRequest,
	}
}

func cors(ctx *routing.Context) error {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", string(ctx.Request.Header.Peek("Origin")))
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "false")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,POST,DELETE")
	ctx.Response.Header.Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)

	if string(ctx.Method()) == "OPTIONS" {
		ctx.Abort()
	}
	if err := ctx.Next(); err != nil {
		if httpError, ok := err.(routing.HTTPError); ok {
			ctx.Response.SetStatusCode(httpError.StatusCode())
		} else {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		}

		b, err := json.Marshal(err)
		if err != nil {
			respond.WithJSON(ctx, fasthttp.StatusInternalServerError, map[string]interface{}{
				"error": err},
			)
			logger.Error("cors marshal", err)
			return nil
		}
		ctx.SetContentType("application/json")
		ctx.SetBody(b)
	}
	return nil
}
