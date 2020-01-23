package server

import (
	"encoding/json"
	"time"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/core/v1"
	"github.com/button-tech/utils-price-tool/core/v2"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	t "github.com/button-tech/utils-price-tool/types"
)

type Core struct {
	S     *fasthttp.Server
	R     *routing.Router
	G     *routing.RouteGroup
	Gv2   *routing.RouteGroup
	store *cache.Cache
}

func New(c *cache.Cache) (core *Core) {
	core = &Core{
		R:     routing.New(),
		store: c,
	}
	core.R.Use(cors)
	core.initBaseRoute()
	core.fs()

	v1.API(core.G, c)
	v2.API(core.Gv2, c)
	return
}

func (core *Core) initBaseRoute() {
	core.G = core.R.Group("/courses/v1")
	core.Gv2 = core.R.Group("/courses/v2")
}

func (core *Core) fs() {
	core.S = &fasthttp.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      core.R.HandleRequest,
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
			respond.WithJSON(ctx, fasthttp.StatusInternalServerError, t.Payload{
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
