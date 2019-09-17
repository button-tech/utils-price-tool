package tasks

import (
	"fmt"
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

func NewGetGroupTask(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)
	wg := sync.WaitGroup{}


	//a, err := slip0044.AddTrustHexToSlip()
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//for _, i := range a {
	//	fmt.Println(i)
	//}

	go func() {
		for ; true; <-ticker {

			// go to compare

			wg.Add(1)
			go func() {
				defer wg.Done()

				res, err := cont.Service.GetCRCPrices()
				if err != nil {
					log.Println(err)
					return
				}

				cont.StoreCRC.Update(res)
			}()

			// go to trust-wallet
			tokens := services.InitRequestData()

			ch := make(chan storetrustwallet.GotPrices, 10)
			var stored []storetrustwallet.GotPrices

			for _, t := range tokens.Tokens {
					got, err := cont.Service.GetPricesCMC(&t)
					if err != nil {
						log.Println(err)
						return
					}
					ch <- got

				item := <-ch
				stored = append(stored, item)
			}

			fmt.Println(runtime.NumGoroutine())
			wg.Wait()
			cont.Store.Update(&stored)
		}
	}()

	cont.Store.Get()
	cont.StoreCRC.Get()
}
