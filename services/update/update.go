package update

import (
	"runtime"
	"sync"
	"time"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/prices"
	"github.com/button-tech/utils-price-tool/pkg/slip0044"
	"github.com/button-tech/utils-price-tool/platforms"
)

type updateWorker func(wg *sync.WaitGroup, prices *prices.PricesData)

func getUpdateWorkers() []updateWorker {
	return []updateWorker{
		platforms.CmcUpdateWorker,
		platforms.HuobiUpdateWorker,
		platforms.TrustUpdateWorker,
		platforms.CrcUpdateWorker,
	}
}

func Start(p *prices.PricesData) {
	converted, err := slip0044.AddTrustHexBySlip()
	if err != nil {
		logger.Error("AddTrustHexBySlip", err)
		return
	}

	var wg sync.WaitGroup
	uw := getUpdateWorkers()
	t := time.NewTicker(time.Minute * 7)
	for ; true; <-t.C {
		start := time.Now()
		if err := p.SetTopList(converted); err != nil {
			logger.Error("GetTopList", err)
			continue
		}

		for _, worker := range uw {
			wg.Add(1)
			go worker(&wg, p)
		}

		logger.Info("Count goroutines: ", runtime.NumGoroutine())
		wg.Wait()

		end := time.Since(start)
		logger.Info("Time EXEC:", end)
	}
}
