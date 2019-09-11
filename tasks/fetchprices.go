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
	End   time.Time
}

func NewGetGroupTask(cont *DuiCont) {
	ticker := time.Tick(cont.TimeOut)

	go func() {
		for range ticker {
			wg := sync.WaitGroup{}

			// go top list
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				list, err := cont.Service.GetTopList()
				if err != nil {
					log.Println(err)
					return
				}

				cont.StoreList.Update(list)

			}(&wg)

			// go to compare
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				res, err := cont.Service.GetCRCPrices()
				if err != nil {
					log.Println(err)
					return
				}

				cont.StoreCRC.Update(res)
			}(&wg)

			// go to trust-wallet
			tokens := services.InitRequestData()

			ch := make(chan storage.GotPrices, 4)
			stored := make([]storage.GotPrices, 0)

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
	cont.StoreList.Get()
}
