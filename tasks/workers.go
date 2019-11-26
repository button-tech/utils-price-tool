package tasks

import (
	"github.com/button-tech/logger"
	"github.com/pkg/errors"
	"sync"

	"github.com/button-tech/utils-price-tool/services"
)

type mappingWorker func(wg *sync.WaitGroup, service *services.Service, store setter)

func cmcWorker(wg *sync.WaitGroup, service *services.Service, store setter) {
	tokens := service.CreateCMCRequestData()

	var tokensWG sync.WaitGroup
	for _, t := range tokens {
		tokensWG.Add(1)
		go func(token services.TokensWithCurrency, tWG *sync.WaitGroup) {
			got, err := service.GetPricesCMC(token)
			if err != nil {
				logger.Error("cmcWorker", err)
			} else {
				store.Set("cmc", got)
			}
			tWG.Done()
		}(t, &tokensWG)
	}
	tokensWG.Wait()
	wg.Done()
}

func crcWorker(wg *sync.WaitGroup, service *services.Service, store setter) {
	res := service.GetPricesCRC()
	if res == nil {
		logger.Error("crcWorker", errors.New("getPricesCRC has nil object"))
		return
	}

	if len(res) > 0 {
		store.Set("crc", res)
	}
	defer wg.Done()
}

func huobiWorker(wg *sync.WaitGroup, service *services.Service, store setter) {
	res, err := service.GetPricesHUOBI()
	if err != nil {
		logger.Error("huobiWorker", err)
	} else {
		store.Set("huobi", res)
	}
	wg.Done()
}
