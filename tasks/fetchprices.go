package tasks

import (
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/slip0044"
	"github.com/button-tech/utils-price-tool/storage"
	"log"
	"runtime"
	"sync"
	"time"
)

type DuiCont struct {
	TimeOut time.Duration
	Service services.Service
	Store   storage.Cached
}

type TickerMeta struct {
	Start time.Time
	End   time.Time
}

//Pool of workers
func NewGetGroup(cont *DuiCont) {
	t := time.NewTicker(cont.TimeOut)

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
		topList, err := cont.Service.GetTopList(converted)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, worker := range workList {
			wg.Add(1)
			go worker(&wg, cont, topList)
		}

		log.Printf("Count goroutines: %v", runtime.NumGoroutine())
		wg.Wait()

		end := time.Since(start)
		log.Println("Time EXEC:", end)
	}
}
