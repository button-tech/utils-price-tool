package tasks

import (
	"runtime"
	"sync"
	"time"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/slip0044"
	"github.com/button-tech/utils-price-tool/services"
)

type worker func(wg *sync.WaitGroup, prices *services.Prices)

func FetchGroup(p *services.Prices) {
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
		if err := p.GetTopList(converted); err != nil {
			logger.Error("GetTopList", err)
			continue
		}

		for _, worker := range ws {
			wg.Add(1)
			go worker(&wg, p)
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

func cmcWorker(wg *sync.WaitGroup, p *services.Prices) {
	defer wg.Done()
	tokens := p.CreateCMCRequestData()

	var tokensWG sync.WaitGroup
	for _, t := range tokens {
		tokensWG.Add(1)
		go func(token services.TokensWithCurrency, tWG *sync.WaitGroup) {
			defer tWG.Done()
			if err := p.SetPricesCMC(token); err != nil {
				logger.Error("cmcWorker", err)
				return
			}
		}(t, &tokensWG)
	}
	tokensWG.Wait()
}

func crcWorker(wg *sync.WaitGroup, p *services.Prices) {
	defer wg.Done()
	p.SetPricesCRC()
}

func huobiWorker(wg *sync.WaitGroup, p *services.Prices) {
	defer wg.Done()
	if err := p.PricesHUOBI(); err != nil {
		logger.Error("huobiWorker", err)
		return
	}
}

func trustV2Worker(wg *sync.WaitGroup, p *services.Prices) {
	defer wg.Done()
	var inWG sync.WaitGroup
	for _, v := range p.TrustV2Coins {
		inWG.Add(1)
		go func(inWg *sync.WaitGroup, price services.PricesTrustV2) {
			defer inWG.Done()
			if err := p.PricesTrustV2(price); err != nil {
				logger.Error("trustV2Worker", err)
				return
			}
		}(&inWG, v)
	}
	inWG.Wait()
}
