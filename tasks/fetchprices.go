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

type Worker func(wg *sync.WaitGroup, cont *DuiCont, list map[string]string)

// CMC worker
func CMCWorker(wg *sync.WaitGroup, cont *DuiCont, list map[string]string) {
	tokens := services.CreateRequestData(list)
	tokensWG := sync.WaitGroup{}

	for _, t := range tokens.Tokens {
		tokensWG.Add(1)
		go func(token services.TokensWithCurrency, tWG *sync.WaitGroup) {
			got, err := cont.Service.GetPricesCMC(token)
			if err != nil {
				log.Println(err)
				return
			}
			cont.Store.Set("cmc", got)

			tWG.Done()
		}(t, &tokensWG)
	}
	tokensWG.Wait()
	wg.Done()
}

// CRC worker
func CRCWorker(wg *sync.WaitGroup, cont *DuiCont, list map[string]string) {
	res, err := cont.Service.GetPricesCRC(list)
	if err != nil {
		log.Println(err)
		return
	}
	cont.Store.Set("crc", res)
	wg.Done()
}

//Pool of workers
func NewGetGroup(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)

	converted, err := slip0044.AddTrustHexBySlip()
	if err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}

	workList := []Worker{
		CMCWorker,
		CRCWorker,
	}

	for ; true; <-ticker {
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
