package api

import (
	"encoding/json"
	"github.com/button-tech/logger"
	"net/http"
	"time"

	"github.com/button-tech/utils-price-tool/pkg/storage"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Core              *fasthttp.Server
	R                 *routing.Router
	G                 *routing.RouteGroup
	Gv2               *routing.RouteGroup
	store             getter
	privateCurrencies map[string][]string
}

type apiController struct {
	privateCurrencies map[string][]string
	store             getter
}

type getter interface {
	Get(a storage.Api) (f storage.FiatMap, err error)
}

func NewServer(store *storage.Cache) *Server {
	server := Server{
		R:                 routing.New(),
		store:             store,
		privateCurrencies: privateCurrencies(),
	}
	server.fs()
	server.R.Use(cors)
	server.initBaseRoute()
	server.initCoursesAPI()
	server.initCoursesAPIv2()

	return &server
}

func (s *Server) fs() {
	s.Core = &fasthttp.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      s.R.HandleRequest,
	}
}

func (s *Server) initBaseRoute() {
	s.G = s.R.Group("/courses/v1")
	s.Gv2 = s.R.Group("/courses/v2")
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
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
		}

		b, err := json.Marshal(err)
		if err != nil {
			respondWithJSON(ctx, fasthttp.StatusInternalServerError, map[string]interface{}{
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

func respondWithJSON(ctx *routing.Context, code int, payload map[string]interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(payload); err != nil {
		logger.Error("write answer", err)
	}
}
