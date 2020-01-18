package v1

import "github.com/button-tech/utils-price-tool/core/internal/handle"

// PricesV1
//
// swagger:response v1Courses
type PricesV1 []handle.Response

// PricesV1 swagger:route POST /v1/prices PricesV1 return prices by input slips tokens
//
// Gets prices by slip tokens
//
// Responses:
//   200: v1Courses
//   400: badRequestResponse

// Error while processing request
//
// swagger:parameters PricesV1

//type inputSlips struct {
//	// in: body
//	Body handle.Data
//}
