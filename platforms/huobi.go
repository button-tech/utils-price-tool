package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"strconv"
	"sync"
)

const urlHuobi = "https://api.hbdm.com/api/v1/contract_index"

func HuobiUpdateWorker(wg *sync.WaitGroup, p *Prices) {
	defer wg.Done()
	if err := p.pricesHuobi(); err != nil {
		logger.Error("huobiWorker", err)
		return
	}
}

func (p *Prices) pricesHuobi() error {
	var huobi huobiResponse

	res, err := req.Get(urlHuobi)
	if err != nil {
		return errors.Wrap(err, "huobi")
	}

	if res.Response().StatusCode != 200 {
		return errors.Wrap(errors.New("error"), "huobi")
	}

	if err := res.ToJSON(&huobi); err != nil {
		return errors.Wrap(err, "toJSON huobi")
	}

	huobiMapping(&huobi, p.List, p.store)

	return nil
}

func huobiMapping(h *huobiResponse, list map[string]string, store *cache.Cache) {
	var wg sync.WaitGroup
	wg.Add(len(h.Data))
	for _, v := range h.Data {
		go func(v huobiData, wg *sync.WaitGroup) {
			if val, ok := list[v.Symbol]; ok {
				defer wg.Done()
				var details cache.Details
				details.Price = strconv.FormatFloat(v.IndexPrice, 'f', -1, 64)
				k := cache.GenKey("huobi", "usd", val)
				store.Set(k, details)
			}
		}(v, &wg)
	}
	wg.Wait()
}
