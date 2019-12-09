package respond

import (
	"encoding/json"
	"fmt"
	"github.com/button-tech/logger"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"io/ioutil"
)

const (
	pathToSwaggerJSON = "./swagger.json"
)

// Payload format: "where: from"
// i.e. "request: json.Unmarshal"
type Error struct {
	API     string
	Func    string
	Err     error
	Payload string
}

func WithJSON(ctx *routing.Context, code int, payload map[string]interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(payload); err != nil {
		logger.Error("write answer", err)
	}
}

func WithWrapErrJSON(ctx *routing.Context, code int, e Error, payload map[string]interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	logger.Error(e.API, e.Err, logger.Params{
		e.Func: e.Payload,
	})
	if payload == nil {
		payload = map[string]interface{}{"error": e.Err.Error()}
	}
	if err := json.NewEncoder(ctx).Encode(payload); err != nil {
		logger.Error("write answer", err)
	}
}

func Payload(where, from string) string {
	return fmt.Sprintf("%s: %s", where, from)
}

func SwaggerJSONHandler(v string) (f func(ctx *routing.Context) error) {
	const funcName = "swaggerJSON"
	f = func(ctx *routing.Context) error {
		plan, err := ioutil.ReadFile(pathToSwaggerJSON)
		if err != nil {
			WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, Error{
				API:     v,
				Func:    funcName,
				Err:     err,
				Payload: "ReadFile",
			}, nil)
			return nil
		}

		var data interface{}
		err = json.Unmarshal(plan, &data)
		if err != nil {
			WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, Error{
				API:     v,
				Func:    funcName,
				Err:     err,
				Payload: err.Error(),
			}, nil)
			return nil
		}
		WithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"data": data})
		return nil
	}
	return
}
