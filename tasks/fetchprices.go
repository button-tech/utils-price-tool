package tasks

import (
	"github.com/utils-price-tool/services"
	"github.com/utils-price-tool/storage"
	"log"
	"sync"
	"time"
)

func NewGetGroupTask(timeout time.Duration, s services.Service, store storage.Storage) {
	ticker := time.Tick(timeout)
	wg := sync.WaitGroup{}

	go func() {
		for _ = range ticker {
			tokens := services.InitRequestData()

			ch := make(chan storage.GotPrices, 3)
			var stored []storage.GotPrices

			for _, t := range tokens.Tokens {
				wg.Add(1)

				go func(wg *sync.WaitGroup, t *services.TokensWithCurrency, ch chan storage.GotPrices) {
					defer wg.Done()

					got, err := s.GetPricesCMC(t)
					if err != nil {
						log.Println(err)
						return
					}

					ch <- got

				}(&wg, &t, ch)
				item := <- ch
				stored = append(stored, item)
			}
			wg.Wait()
			store.Update(&stored)

		}
	}()

	store.Get()
}
