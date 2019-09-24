package tasks

import (
	"fmt"
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

type Worker func(wg *sync.WaitGroup, cont *DuiCont)

func Worker2(wg *sync.WaitGroup, cont *DuiCont) {
	tokens := services.InitRequestData()
	for _, t := range tokens.Tokens {

		got, err := cont.Service.GetPricesCMC(t)
		if err != nil {
			log.Println(err)
			return
		}
		cont.Store.Set("cmc", got)
	}
	wg.Done()
}

func Worker1(wg *sync.WaitGroup, cont *DuiCont) {
	res, err := cont.Service.GetPricesCRC()
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
		Worker1,
		Worker2,
	}

	for ; true; <-ticker {
		start := time.Now()
		topList, err := cont.Service.GetTopList(converted)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(len(topList))

		for _, worker := range workList {
			wg.Add(1)
			go worker(&wg, cont)
		}
		log.Printf("Count goroutines: %v", runtime.NumGoroutine())
		wg.Wait()

		end := time.Since(start)
		log.Println("Time EXEC:", end)
	}
}
