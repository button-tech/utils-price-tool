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
				return
			}
			store.Set("cmc", got)
			defer tWG.Done()
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
		return
	}
	store.Set("huobi", res)
	defer wg.Done()
}

//func trustV2Worker(wg *sync.WaitGroup, service services.Service, store setter) {
//	var inWG sync.WaitGroup
//	for _, v := range service.TrustV2Coins {
//		inWG.Add(1)
//		go func() {
//
//		}()
//	}
//}
