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
	TimeOut   time.Duration
	Service   services.Service
	Store     storage.Cached
}

type TickerMeta struct {
	Start time.Time
	End   time.Time
}

//Pool of workers
func NewGetGroupTask(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)

	converted, err := slip0044.AddTrustHexBySlip()
	if err != nil {
		log.Println(err)
		return
	}
	wg := sync.WaitGroup{}


	for ; true; <-ticker {
		start := time.Now()
		topList, err := cont.Service.GetTopList(converted)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(len(topList))

		//cryptoForCRC := make([]string, 0)
		//for _, v := range topList {
		//	cryptoForCRC = append(cryptoForCRC, v)
		//}
		wg.Add(2)

		// go to compare
		go func(wg *sync.WaitGroup) {
			res, err := cont.Service.GetPricesCRC()
			if err != nil {
				log.Println(err)
				return
			}
			cont.Store.Set("crc", res)
			wg.Done()
		}(&wg)


		// go to trust-wallet
		go func(wg *sync.WaitGroup) {
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
		}(&wg)

		wg.Wait()

		log.Printf("Count goroutines: %v", runtime.NumGoroutine())
		end := time.Since(start)
		log.Println("Time EXEC:", end)
	}
}
