package api

import (
	"encoding/json"
	"errors"
	"github.com/button-tech/utils-price-tool/storage"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strings"
)

type request struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

type response struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

type listApi struct {
	API []api `json:"api"`
	//Time             struct {
	//	Start int `json:"start"`
	//	End   int `json:"end"`
	//} `json:"time"`
}

type api struct {
	Name             string   `json:"name"`
	SupportedChanges []string `json:"supported_changes"`
}

func (ac *apiController) getCourses(ctx *routing.Context) error {
	var req request
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		//respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
		//	"error": err.Error(),
		//})
		return err
	}

	a := req.API; switch a {
	case "cmc", "crc":
		result, err := ac.converter(&req, a)
		if err != nil {
			//respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			//	"error":"no matches API changes",
			//})
			return err
		}
		respondWithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{
			"data":result,
		})
		return nil

	default:
		supportedCRC := []string{"0", "1", "24"}
		crc := api{
			Name:             "crc",
			SupportedChanges: supportedCRC,
		}

		supportedCMC := []string{"24"}
		cmc := api{
			Name:             "cmc",
			SupportedChanges: supportedCMC,
		}

		API := []api{crc, cmc}
		list := listApi{API: API}

		respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			"error": &list,
			"description": "please, use these API",
		})
		return nil
	}
}

func (ac *apiController) apiInfo(ctx *routing.Context) error {
	supportedCRC := []string{"0", "1", "24"}
	crc := api{
		Name:             "crc",
		SupportedChanges: supportedCRC,
	}

	supportedCMC := []string{"24"}
	cmc := api{
		Name:             "cmc",
		SupportedChanges: supportedCMC,
	}

	API := []api{crc, cmc}
	list := listApi{API: API}
	respondWithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{
		"api": &list,
	})
	return nil
}

func (ac *apiController) mapping(req *request, api string) []*response {
	result := make([]*response, 0)
	stored := ac.store.Get()[storage.Api(api)]

	for _, c := range req.Currencies {
		price := response{}

		if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
			price.Currency = c

			for _, t := range req.Tokens {
				if val, ok := fiatVal[storage.CryptoCurrency(strings.ToLower(t))]; ok {
					contract := map[string]string{t: val.Price}
					if contract = changesControl(contract, val, req.Change); contract == nil {
						return nil
					} else {
						price.Rates = append(price.Rates, contract)
					}
				}
			}
		}
		if price.Currency != "" {
			result = append(result, &price)
		}
	}
	return result
}

func changesControl(m map[string]string, s *storage.Details, c string) map[string]string {
	switch c {
	case "1":
		if s.ChangePCTHour != "" {
			m["percent_change"] = s.ChangePCTHour

			return m
		}
		return nil
	case "24":
		if s.ChangePCT24Hour != "" {
			m["percent_change"] = s.ChangePCT24Hour
			return m
		}
		return nil
	default:
		return m
	}
}

func (ac *apiController) converter(req *request, api string) ([]*response, error) {
	a := api
	switch a {
	case "cmc", "crc":
		resp := ac.mapping(req, a)
		if resp == nil {
			return nil, errors.New("no matches API")
		}
		return resp, nil
	default:
		return nil, errors.New("no matches API")
	}
}

func (s *Server) initCoursesAPI()  {
	s.G.Post("/prices", s.ac.getCourses)
	s.G.Get("/list", s.ac.apiInfo)
}
