package respond

import (
	"encoding/json"
	"fmt"
	"github.com/button-tech/logger"
	routing "github.com/qiangxue/fasthttp-routing"
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
