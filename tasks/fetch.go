package tasks

import (
	"runtime"
	"sync"
	"time"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/slip0044"
	"github.com/button-tech/utils-price-tool/services"
)

type worker func(wg *sync.WaitGroup, service *services.GetPrices)

func FetchGroup(service *services.GetPrices) {
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
			go worker(&wg, service)
		}

		logger.Info("Count goroutines: ", runtime.NumGoroutine())
		wg.Wait()

		end := time.Since(start)
		logger.Info("Time EXEC:", end)
	}
}

func workers() []worker {
	return []worker{
		cmcWorker,
		crcWorker,
		huobiWorker,
		trustV2Worker,
	}
}

func cmcWorker(wg *sync.WaitGroup, service *services.GetPrices) {
	tokens := service.CreateCMCRequestData()

	var tokensWG sync.WaitGroup
	for _, t := range tokens {
		tokensWG.Add(1)
		go func(token services.TokensWithCurrency, tWG *sync.WaitGroup) {
			if err := service.GetPricesCMC(token); err != nil {
				logger.Error("cmcWorker", err)
				return
			}
			defer tWG.Done()
		}(t, &tokensWG)
	}

	tokensWG.Wait()
	wg.Done()
}

func crcWorker(wg *sync.WaitGroup, service *services.GetPrices) {
	service.GetPricesCRC()
	defer wg.Done()
}

func huobiWorker(wg *sync.WaitGroup, service *services.GetPrices) {
	if err := service.GetPricesHUOBI(); err != nil {
		logger.Error("huobiWorker", err)
		return
	}
	defer wg.Done()
}

func trustV2Worker(wg *sync.WaitGroup, service *services.GetPrices) {
	var inWG sync.WaitGroup
	for _, v := range service.TrustV2Coins {
		inWG.Add(1)
		go func(inWg *sync.WaitGroup, price services.PricesTrustV2) {
			if err := service.GetPricesTrustV2(price); err != nil {
				logger.Error("trustV2Worker", err)
				return
			}
			defer inWG.Done()
		}(&inWG, v)
	}
	inWG.Wait()
	wg.Done()
}
