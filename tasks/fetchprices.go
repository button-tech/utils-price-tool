package tasks

import (
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"log"
	"sync"
	"time"
)

type DuiCont struct {
	TimeOut   time.Duration
	Service   services.Service
	Store     storage.Storage
	StoreList storetoplist.Storage
	StoreCRC  storecrc.Storage
}

type TickerMeta struct {
	Start time.Time
	End time.Time
}

func NewGetGroupTask(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)
	wg := sync.WaitGroup{}

	go func() {
		for _ = range ticker {

			// go to compare
			go func() {
				wg.Add(1)
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

			ch := make(chan storage.GotPrices, 4)
			var stored []storage.GotPrices

			for _, t := range tokens.Tokens {
				wg.Add(1)

				go func(wg *sync.WaitGroup, t *services.TokensWithCurrency, ch chan storage.GotPrices) {
					defer wg.Done()

					got, err := cont.Service.GetPricesCMC(t)
					if err != nil {
						log.Println(err)
						return
					}

					ch <- got

				}(&wg, &t, ch)
				item := <-ch
				stored = append(stored, item)
			}

			wg.Wait()
			cont.Store.Update(&stored)
		}
	}()

	cont.Store.Get()
	cont.StoreCRC.Get()
}
