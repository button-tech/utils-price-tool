package tasks

import (
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"log"
	"runtime"
	"sync"
	"time"
)

type DuiCont struct {
	TimeOut   time.Duration
	Service   services.Service
	Store     storage.Cached
	StoreList storetoplist.Storage
	StoreCRC  storecrc.Storage
}

type TickerMeta struct {
	Start time.Time
	End   time.Time
}

// Pool of workers
func NewGetGroupTask(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)

	//converted, err := slip0044.AddTrustHexBySlip()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	wg := sync.WaitGroup{}

	go func() {
		for ; true; <-ticker {

			//topList, err := cont.Service.GetTopList(converted)
			//if err != nil {
			//	log.Println(err)
			//}
			//
			//cryptoForCRC := make([]string, 0)
			//for _, v := range topList {
			//	cryptoForCRC = append(cryptoForCRC, v)
			//}


			// go to compare
			//wg.Add(1)
			//go func(wg *sync.WaitGroup) {
			//	res, err := cont.Service.GetCRCPrices(cryptoForCRC)
			//	if err != nil {
			//		log.Println(err)
			//		return
			//	}
			//
			//	for k, v := range res {
			//		cont.StoreCRC.Update(k, v)
			//	}
			//
			//	wg.Done()
			//}(&wg)

			// go to trust-wallet
			tokens := services.InitRequestData()

			//ch := make(chan *storetrustwallet.GotPrices, 10)
			//stored := make([]*storetrustwallet.GotPrices, 0)

			for _, t := range tokens.Tokens {
				wg.Add(1)

				go func(t services.TokensWithCurrency, wg *sync.WaitGroup) {
					got, err := cont.Service.GetPricesCMC(t)
					if err != nil {
						log.Println(err)
						return
					}

					cont.Store.Set("cmc", got)
					wg.Done()
				}(t, &wg)

				//item := <-ch
				//stored = append(stored, item)
			}

			log.Printf("Count goroutines: %v", runtime.NumGoroutine())
			wg.Wait()
			//cont.Store.Update(stored)
		}
	}()

	cont.Store.Get()
	//cont.StoreCRC.Get()
}
