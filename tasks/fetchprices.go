package tasks

import (
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"github.com/button-tech/utils-price-tool/storage/storetrustwallet"
	"log"
	"runtime"
	"sync"
	"time"
)

type DuiCont struct {
	TimeOut   time.Duration
	Service   services.Service
	Store     storetrustwallet.Storage
	StoreList storetoplist.Storage
	StoreCRC  storecrc.Storage
}

type TickerMeta struct {
	Start time.Time
	End   time.Time
}

// pool of workers
func NewGetGroupTask(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)
	wg := sync.WaitGroup{}

	//topList, err := cont.Service.GetTopList()
	//if err != nil {
	//	log.Println(err)
	//}


	go func() {
		for ; true; <-ticker {

			// go to compare
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				res, err := cont.Service.GetCRCPrices()
				if err != nil {
					log.Println(err)
					return
				}

				cont.StoreCRC.Update(res)
				wg.Done()
			}(&wg)

			// go to trust-wallet
			tokens := services.InitRequestData()

			ch := make(chan *storetrustwallet.GotPrices, 10)
			stored := make([]*storetrustwallet.GotPrices, 0)

			for _, t := range tokens.Tokens {
				wg.Add(1)

				go func(t services.TokensWithCurrency, wg *sync.WaitGroup) {
					got, err := cont.Service.GetPricesCMC(t)
					if err != nil {
						log.Println(err)
						return
					}

					ch <- got
					wg.Done()
				}(t, &wg)

				item := <-ch
				stored = append(stored, item)
			}

			log.Printf("Count goroutines: %v", runtime.NumGoroutine())
			wg.Wait()
			cont.Store.Update(stored)
		}
	}()

	cont.Store.Get()
	cont.StoreCRC.Get()
}
