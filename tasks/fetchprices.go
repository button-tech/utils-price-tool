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
	wg := sync.WaitGroup{}

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

//func NewGetGroupTask(cont *DuiCont) {
//	ticker := time.Tick(cont.TimeOut)
//	wg := sync.WaitGroup{}
//
//	go func() {
//		for ; true; <-ticker {
//			wg.Add(1)
//			go func(cont *DuiCont, wg *sync.WaitGroup) {
//
//				// go top list
//				wg.Add(1)
//				go func(wg *sync.WaitGroup, cont *DuiCont) {
//					defer wg.Done()
//					list, err := cont.Service.GetTopList()
//					if err != nil {
//						log.Println(err)
//						return
//					}
//
//					cont.StoreList.Update(list)
//
//				}(wg, cont)
//
//				// go to compare
//				wg.Add(1)
//				go func(wg *sync.WaitGroup, cont *DuiCont) {
//					defer wg.Done()
//
//					res, err := cont.Service.GetCRCPrices()
//					if err != nil {
//						log.Println(err)
//						return
//					}
//
//					cont.StoreCRC.Update(res)
//				}(wg, cont)
//
//				// go to trust-wallet
//				tokens := services.InitRequestData()
//
//				ch := make(chan storage.GotPrices, 4)
//				stored := make([]storage.GotPrices, 0)
//
//				go func(wg *sync.WaitGroup, ch chan storage.GotPrices, cont *DuiCont) {
//					wg.Add(1)
//
//					for _, t := range tokens.Tokens {
//						wg.Add(1)
//
//						go func(wg *sync.WaitGroup, t *services.TokensWithCurrency, ch chan storage.GotPrices) {
//							defer wg.Done()
//
//							got, err := cont.Service.GetPricesCMC(t)
//							if err != nil {
//								log.Println(err)
//								return
//							}
//
//							ch <- got
//
//						}(wg, &t, ch)
//						item := <-ch
//						stored = append(stored, item)
//					}
//				}(wg, ch, cont)
//
//				wg.Wait()
//				cont.Store.Update(&stored)
//			}(cont, &wg)
//
//		}
//	}()
//
//	cont.Store.Get()
//	cont.StoreCRC.Get()
//
//}