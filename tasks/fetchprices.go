package tasks

import (
	"github.com/utils-tool_prices/services"
	"github.com/utils-tool_prices/storage"
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
			ch := make(chan storage.Prices, 3)
			var stored []storage.Prices

			for _, t := range tokens.Tokens {
				wg.Add(1)

				go func(wg *sync.WaitGroup, t *services.TokensWithCurrency, ch chan storage.Prices) {
					defer wg.Done()

					got, err := s.GetPricesCMC(t)
					if err != nil {
						log.Println(err)
						return
					}

					var storeItem storage.Prices
					for _, i := range got.Docs {
						storeItem.Currency = got.Currency
						contract := map[string]string{i.Contract: i.Price}
						storeItem.Rates = append(storeItem.Rates, &contract)
					}

					ch <- storeItem

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
