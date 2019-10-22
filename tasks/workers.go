package tasks

import (
	"github.com/button-tech/utils-price-tool/services"
	"log"
	"sync"
)

type mappingWorker func(wg *sync.WaitGroup, service *services.Service, store setter)

func cmcWorker(wg *sync.WaitGroup, service *services.Service, store setter) {
	tokens := service.CreateCMCRequestData()
	tokensWG := sync.WaitGroup{}

	for _, t := range tokens.Tokens {
		tokensWG.Add(1)
		go func(token services.TokensWithCurrency, tWG *sync.WaitGroup) {
			got, err := service.GetPricesCMC(token)
			if err != nil {
				log.Println(err)
				return
			}
			store.Set("cmc", got)

			tWG.Done()
		}(t, &tokensWG)
	}
	tokensWG.Wait()
	wg.Done()
}

func crcWorker(wg *sync.WaitGroup, service *services.Service, store setter) {
	res := service.GetPricesCRC()
	store.Set("crc", res)
	wg.Done()
}

func huobiWorker(wg *sync.WaitGroup, service *services.Service, store setter) {
	res := service.GetPricesHUOBI()
	store.Set("huobi", res)
	wg.Done()
}
