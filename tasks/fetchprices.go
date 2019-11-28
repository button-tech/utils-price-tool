package tasks

import (
	"github.com/button-tech/logger"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/slip0044"
	"github.com/button-tech/utils-price-tool/storage"
)

type DuiCont struct {
	TimeOut time.Duration
	Service *services.Service
	Store   *storage.Cache
}

type TickerMeta struct {
	Start time.Time
	End   time.Time
}

type setter interface {
	Set(a storage.Api, f storage.FiatMap)
}

//Pool of workers
func NewGetGroup(service *services.Service, store setter) {
	t := time.NewTicker(time.Minute * 7)

	converted, err := slip0044.AddTrustHexBySlip()
	if err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}

	workList := []mappingWorker{
		cmcWorker,
		crcWorker,
		huobiWorker,
	}

	for ; true; <-t.C {
		start := time.Now()
		if err := service.GetTopList(converted); err != nil {
			logger.Error("GetTopList", err)
			continue
		}

		for _, worker := range workList {
			wg.Add(1)
			go worker(&wg, service, store)
		}

		logger.Info("Count goroutines: ", runtime.NumGoroutine())
		wg.Wait()

		end := time.Since(start)
		logger.Info("Time EXEC:", end)
	}
}
