package v2

import (
	"encoding/json"

	"github.com/button-tech/utils-price-tool/core/internal/handle"
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/button-tech/utils-price-tool/services"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const (
	v2 = "v2"
)

func (c *controller) courses(ctx *routing.Context) error {
	const funcName = "courses"
	var r handle.Data
	if err := json.Unmarshal(ctx.PostBody(), &r); err != nil {
		respond.WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, respond.Error{
			API:     v2,
			Func:    funcName,
			Err:     err,
			Payload: respond.Payload("request", "json.Unmarshal"),
		}, nil)
		return nil
	}

	unique := handle.Unify(&r)
	resp, err := handle.Reply(&unique, v2, c.store, nil)
	if err != nil {
		respond.WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, respond.Error{
			API:     v2,
			Func:    funcName,
			Err:     err,
			Payload: respond.Payload("response", "handle.Reply"),
		}, map[string]interface{}{"api": supportInfo(), "error": "please, use these API"})
		return nil
	}

	respond.WithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"data": resp})
	return nil
}

func (c *controller) info(ctx *routing.Context) error {
	respond.WithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"api": supportInfo()})
	return nil
}

func supportInfo() []handle.APIs {
	var apis []handle.APIs
	supportedNewTrust := []string{"0", "24"}
	supportedPureCMC := []string{"0", "24", "7d"}
	apis = append(apis, handle.APIs{
		Name:             "ntrust",
		SupportedChanges: supportedNewTrust,
		SupportedFiats:   services.TrustV2Coins,
	}, handle.APIs{
		Name:             "pcmc",
		SupportedChanges: supportedPureCMC,
		SupportedFiats:   services.PureCMCCoins,
	})
	return apis
}
