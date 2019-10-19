package api

import (
	"encoding/json"
	"github.com/button-tech/utils-price-tool/storage"
	routing "github.com/qiangxue/fasthttp-routing"
	"log"
)

type Server struct {
	R     *routing.Router
	G     *routing.RouteGroup
	ac    *apiController
	store storage.Cached
}

func NewServer(store storage.Cached) *Server {
	server := Server{
		R:     routing.New(),
		store: store,
	}
	server.initBaseRoute()
	server.initCoursesAPI()

	return &server
}

func (s *Server) initBaseRoute() {
	s.G = s.R.Group("/courses/v1")
	s.ac = &apiController{store: s.store}
}

type apiController struct {
	store storage.Cached
}

func respondWithJSON(ctx *routing.Context, code int, payload map[string]interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(payload); err != nil {
		log.Println("write answer", err)
	}
}