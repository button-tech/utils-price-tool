package tasks

import (
	"runtime"
	"sync"
	"time"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/slip0044"
	"github.com/button-tech/utils-price-tool/pkg/storage"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/pkg/errors"
)

type setter interface {
	Set(a storage.Api, f storage.FiatMap)
}

type worker func(wg *sync.WaitGroup, service *services.Service, store setter)

func FetchGroup(service *services.Service, store setter) {
	converted, err := slip0044.AddTrustHexBySlip()
	if err != nil {
		logger.Error("AddTrustHexBySlip", err)
		return
	}

	var wg sync.WaitGroup
	ws := workers()
	t := time.NewTicker(time.Minute * 7)
	for ; true; <-t.C {
		start := time.Now()
		if err := service.GetTopList(converted); err != nil {
			logger.Error("GetTopList", err)
			continue
		}

		for _, worker := range ws {
			wg.Add(1)
			go worker(&wg, service, store)
		}

		logger.Info("Count goroutines: ", runtime.NumGoroutine())
		wg.Wait()

		end := time.Since(start)
		logger.Info("Time EXEC:", end)
	}
}

func workers() []worker {
	return []worker{cmcWorker,
		crcWorker,
		huobiWorker,
		trustV2Worker,
	}
}

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

func trustV2Worker(wg *sync.WaitGroup, service *services.Service, store setter) {
	var inWG sync.WaitGroup
	for _, v := range service.TrustV2Coins {
		inWG.Add(1)
		go func(inWg *sync.WaitGroup, price services.PricesTrustV2) {
			got, err := service.GetPricesTrustV2(price)
			if err != nil {
				logger.Error("trustV2Worker", err)
				return
			}
			store.Set("ntrust", got)
			defer inWG.Done()
		}(&inWG, v)
	}
	inWG.Wait()
	wg.Done()
}
